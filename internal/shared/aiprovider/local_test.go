package aiprovider

import (
	"context"
	"testing"

	"github.com/JarvanDante/my_service/internal/model/entity"
	"github.com/JarvanDante/my_service/internal/shared/aiprotocol"
)

func TestBuildFaceSwapJobFromURLs(t *testing.T) {
	job, err := BuildFaceSwapJob(context.Background(), SubmitInput{
		TaskNo:   "A1",
		BizType:  entity.AiBizFaceSwap,
		InputURL: "http://127.0.0.1:19000/my-storage/my/image/a/face.bnc",
		Params: map[string]any{
			"target_url": "http://127.0.0.1:19000/my-storage/my/image/b/tpl.bnc",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if job.JobID != "A1" || job.MediaType != aiprotocol.MediaPhoto {
		t.Fatalf("%+v", job)
	}
	if job.Source.Key != "my/image/a/face.bnc" || job.Target.Key != "my/image/b/tpl.bnc" {
		t.Fatalf("source=%+v target=%+v", job.Source, job.Target)
	}
	if job.Output.Prefix != "my/ai/out/A1/" || job.Output.Bucket != aiprotocol.DefaultBucket {
		t.Fatalf("%+v", job.Output)
	}
}

func TestBuildFaceSwapJobRejectsVideo(t *testing.T) {
	_, err := BuildFaceSwapJob(context.Background(), SubmitInput{
		TaskNo:  "A1",
		BizType: entity.AiBizFaceSwap,
		Params:  map[string]any{"media_type": "video", "source_url": "a.jpg", "target_url": "b.jpg"},
	})
	if err == nil {
		t.Fatal("want video rejected")
	}
}

func TestBuildFaceSwapJobRequiresTarget(t *testing.T) {
	_, err := BuildFaceSwapJob(context.Background(), SubmitInput{
		TaskNo:   "A1",
		BizType:  entity.AiBizFaceSwap,
		InputURL: "my/image/a/face.bnc",
	})
	if err == nil {
		t.Fatal("want target required")
	}
}

func TestBuildUndressJobFromURL(t *testing.T) {
	job, err := BuildUndressJob(context.Background(), SubmitInput{
		TaskNo:   "A2",
		BizType:  entity.AiBizUndress,
		InputURL: "http://127.0.0.1:19000/my-storage/my/image/c/photo.bnc",
	})
	if err != nil {
		t.Fatal(err)
	}
	if job.JobID != "A2" || job.Biz != aiprotocol.BizUndress {
		t.Fatalf("%+v", job)
	}
	if job.Source.Key != "my/image/c/photo.bnc" || job.Target.Key != "" {
		t.Fatalf("source=%+v target=%+v", job.Source, job.Target)
	}
	if job.Output.Prefix != "my/ai/out/A2/" {
		t.Fatalf("%+v", job.Output)
	}
}

func TestBuildUndressJobRequiresSource(t *testing.T) {
	_, err := BuildUndressJob(context.Background(), SubmitInput{
		TaskNo:  "A2",
		BizType: entity.AiBizUndress,
	})
	if err == nil {
		t.Fatal("want source required")
	}
}

func TestLocalSubmitRejectsOtherBiz(t *testing.T) {
	p := &localProvider{}
	_, err := p.Submit(context.Background(), SubmitInput{
		TaskNo:  "A3",
		BizType: entity.AiBizTextToImage,
	})
	if err == nil {
		t.Fatal("want other biz rejected")
	}
}
