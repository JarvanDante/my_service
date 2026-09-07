package logic

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/pquerna/otp/totp"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/modules/admin/domain"
	"github.com/JarvanDante/my_service/internal/modules/admin/service"
)

func TestLoginTotpFlow(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	salt := "testsalt"
	repo := newFakeAdminRepo(&entity.AdminUser{
		Id: 1, Username: "admin", Password: gmd5.MustEncryptString("admin123" + salt),
		Salt: salt, Nickname: "超管", RoleId: 1, Status: 1,
	})
	pending := newMemTotpPending()
	svc := newAdmin(repo, pending)

	oldIssue := issueAdminToken
	issueAdminToken = func(_ context.Context, id int64) (string, error) { return "tok-1", nil }
	t.Cleanup(func() { issueAdminToken = oldIssue })

	if _, err := svc.Login(ctx, service.LoginInput{Username: "admin", Password: "wrong"}); err == nil {
		t.Fatal("wrong password should fail")
	}

	first, err := svc.Login(ctx, service.LoginInput{Username: "admin", Password: "admin123"})
	if err != nil {
		t.Fatal(err)
	}
	if !first.NeedTotp || first.TotpBound || first.Token != "" || first.TotpSecret == "" || !strings.HasPrefix(first.TotpQR, "data:image/png;base64,") {
		t.Fatalf("unbound first step: %+v", first)
	}

	again, err := svc.Login(ctx, service.LoginInput{Username: "admin", Password: "admin123"})
	if err != nil {
		t.Fatal(err)
	}
	if again.TotpSecret != first.TotpSecret {
		t.Fatal("pending secret should be reused")
	}

	if _, err = svc.Login(ctx, service.LoginInput{Username: "admin", Password: "admin123", TotpCode: "000000"}); err == nil {
		t.Fatal("wrong totp should fail")
	}
	if repo.admins[1].TotpSecret != "" {
		t.Fatal("wrong totp must not bind")
	}

	code, err := totp.GenerateCode(first.TotpSecret, now)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := svc.Login(ctx, service.LoginInput{Username: "admin", Password: "admin123", TotpCode: code, Ip: "1.1.1.1"})
	if err != nil {
		t.Fatal(err)
	}
	if ok.Token != "tok-1" || ok.NeedTotp || ok.Admin == nil || ok.Admin.Username != "admin" {
		t.Fatalf("bind+login: %+v", ok)
	}
	if !totpBound(repo.admins[1]) || repo.lastIP != "1.1.1.1" {
		t.Fatal("should persist bind and login ip")
	}
	if secret, _ := pending.Get(ctx, 1); secret != "" {
		t.Fatal("pending secret should be cleared after bind")
	}

	step, err := svc.Login(ctx, service.LoginInput{Username: "admin", Password: "admin123"})
	if err != nil {
		t.Fatal(err)
	}
	if !step.NeedTotp || !step.TotpBound || step.TotpQR != "" || step.TotpSecret != "" {
		t.Fatalf("bound challenge should hide secret: %+v", step)
	}

	code2, err := totp.GenerateCode(repo.admins[1].TotpSecret, now)
	if err != nil {
		t.Fatal(err)
	}
	done, err := svc.Login(ctx, service.LoginInput{Username: "admin", Password: "admin123", TotpCode: code2})
	if err != nil {
		t.Fatal(err)
	}
	if done.Token != "tok-1" {
		t.Fatalf("bound login token=%q", done.Token)
	}

	if err = svc.ResetTotp(ctx, 1); err != nil {
		t.Fatal(err)
	}
	if totpBound(repo.admins[1]) {
		t.Fatal("reset should unbind")
	}
}

func TestLoginDisabledAndMissing(t *testing.T) {
	ctx := context.Background()
	salt := "testsalt"
	repo := newFakeAdminRepo(&entity.AdminUser{
		Id: 2, Username: "off", Password: gmd5.MustEncryptString("admin123" + salt), Salt: salt, Status: 0,
	})
	svc := newAdmin(repo, newMemTotpPending())
	if _, err := svc.Login(ctx, service.LoginInput{Username: "off", Password: "admin123"}); err == nil {
		t.Fatal("disabled account should fail")
	}
	if _, err := svc.Login(ctx, service.LoginInput{Username: "ghost", Password: "admin123"}); err == nil {
		t.Fatal("missing account should fail")
	}
}

type fakeAdminRepo struct {
	stubRepo
	admins map[int64]*entity.AdminUser
	lastIP string
}

func newFakeAdminRepo(list ...*entity.AdminUser) *fakeAdminRepo {
	m := make(map[int64]*entity.AdminUser)
	for _, a := range list {
		cp := *a
		m[a.Id] = &cp
	}
	return &fakeAdminRepo{admins: m}
}

func (f *fakeAdminRepo) FindByUsername(_ context.Context, username string) (*entity.AdminUser, error) {
	for _, a := range f.admins {
		if a.Username == username {
			cp := *a
			return &cp, nil
		}
	}
	return nil, nil
}

func (f *fakeAdminRepo) FindById(_ context.Context, id int64) (*entity.AdminUser, error) {
	a := f.admins[id]
	if a == nil {
		return nil, nil
	}
	cp := *a
	return &cp, nil
}

func (f *fakeAdminRepo) UpdateLoginInfo(_ context.Context, id int64, ip string) error {
	f.lastIP = ip
	if a := f.admins[id]; a != nil {
		a.LastLoginAt = gtime.Now()
		a.LastIp = ip
	}
	return nil
}

func (f *fakeAdminRepo) BindTotp(_ context.Context, id int64, secret string) error {
	a := f.admins[id]
	if a == nil {
		return errors.New("missing")
	}
	a.TotpSecret = secret
	a.TotpBoundAt = gtime.Now()
	return nil
}

func (f *fakeAdminRepo) ResetTotp(_ context.Context, id int64) error {
	a := f.admins[id]
	if a == nil {
		return errors.New("missing")
	}
	a.TotpSecret = ""
	a.TotpBoundAt = nil
	return nil
}

type stubRepo struct{}

func (stubRepo) ListRoles(context.Context) ([]*entity.AdminRole, error) { return nil, errStub }
func (stubRepo) FindRoleById(context.Context, int64) (*entity.AdminRole, error) {
	return nil, errStub
}
func (stubRepo) FindRoleByCode(context.Context, string) (*entity.AdminRole, error) {
	return nil, errStub
}
func (stubRepo) CreateRole(context.Context, *entity.AdminRole) (int64, error) { return 0, errStub }
func (stubRepo) UpdateRole(context.Context, int64, string, string, int, string) error {
	return errStub
}
func (stubRepo) DeleteRole(context.Context, int64) error                 { return errStub }
func (stubRepo) CountAdminsByRoleId(context.Context, int64) (int, error) { return 0, errStub }
func (stubRepo) ListPermissions(context.Context) ([]*entity.AdminPermission, error) {
	return nil, errStub
}
func (stubRepo) FindPermissionsByIds(context.Context, []int64) ([]*entity.AdminPermission, error) {
	return nil, errStub
}
func (stubRepo) FindPermissionById(context.Context, int64) (*entity.AdminPermission, error) {
	return nil, errStub
}
func (stubRepo) CreatePermission(context.Context, *entity.AdminPermission) (int64, error) {
	return 0, errStub
}
func (stubRepo) UpdatePermission(context.Context, *entity.AdminPermission) error { return errStub }
func (stubRepo) DeletePermission(context.Context, int64) error                   { return errStub }
func (stubRepo) CountPermissionChildren(context.Context, int64) (int, error)     { return 0, errStub }
func (stubRepo) ListAdmins(context.Context, int, int) ([]*entity.AdminUser, int, error) {
	return nil, 0, errStub
}
func (stubRepo) CreateAdmin(context.Context, *entity.AdminUser) (int64, error) { return 0, errStub }
func (stubRepo) UpdateAdmin(context.Context, int64, string, int64, int, string, string) error {
	return errStub
}
func (stubRepo) DeleteAdmin(context.Context, int64) error { return errStub }

var (
	errStub                   = errors.New("stub")
	_       domain.Repository = (*fakeAdminRepo)(nil)
	_       service.IAdmin    = (*sAdmin)(nil)
)
