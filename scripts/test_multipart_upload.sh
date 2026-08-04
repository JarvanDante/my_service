#!/usr/bin/env bash
# 端到端测试 MinIO 分片上传:
#   login → init → parts → presign → PUT 各片 → complete
#
# 用法:
#   ./scripts/test_multipart_upload.sh                          # 默认用 oceans.mp4
#   ./scripts/test_multipart_upload.sh /path/to/other.mp4
#   PART_SIZE=5242880 ./scripts/test_multipart_upload.sh        # 强制 5MiB 多分片
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8000}"
USER="${ADMIN_USER:-admin}"
PASS="${ADMIN_PASS:-admin123}"
# 默认测试片: ~/Downloads/oceans.mp4 (~22MB)
DEFAULT_FILE="/Users/wangdante/Downloads/oceans.mp4"
FILE="${1:-$DEFAULT_FILE}"
PART_SIZE="${PART_SIZE:-8388608}"

if [[ ! -f "$FILE" ]]; then
  echo "文件不存在: $FILE"
  echo "用法: $0 [video.mp4]"
  exit 1
fi

SIZE=$(wc -c < "$FILE" | tr -d ' ')
NAME=$(basename "$FILE")
TMPDIR="${TMPDIR:-/tmp}/mp_upload_$$"
mkdir -p "$TMPDIR"
trap 'rm -rf "$TMPDIR"' EXIT

echo "==> 1) 登录 $BASE_URL"
curl -sS -X POST "$BASE_URL/backend/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"username\":\"$USER\",\"password\":\"$PASS\"}" > "$TMPDIR/login.json"
TOKEN=$(python3 -c 'import json,sys; d=json.load(open(sys.argv[1])); assert d.get("code")==0,d; print(d["data"]["token"])' "$TMPDIR/login.json")
echo "    ok"

echo "==> 2) init (file=$NAME size=$SIZE part_size=$PART_SIZE)"
curl -sS -X POST "$BASE_URL/backend/media/multipart/init" \
  -H "Authorization: $TOKEN" -H 'Content-Type: application/json' \
  -d "{\"filename\":\"$NAME\",\"purpose\":\"video\",\"content_type\":\"video/mp4\",\"size\":$SIZE,\"part_size\":$PART_SIZE}" \
  > "$TMPDIR/init.json"
python3 - "$TMPDIR/init.json" <<'PY'
import json,sys
d=json.load(open(sys.argv[1]))
assert d.get("code")==0, d
print("    upload_id =", d["data"]["upload_id"])
print("    object_key=", d["data"]["object_key"])
print("    part_count=", d["data"]["part_count"])
json.dump(d["data"], open(sys.argv[1]+".data","w"))
PY
UPLOAD_ID=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["upload_id"])' "$TMPDIR/init.json.data")
PART_COUNT=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["part_count"])' "$TMPDIR/init.json.data")
PART_SIZE=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["part_size"])' "$TMPDIR/init.json.data")

echo "==> 3) parts (续传探测)"
curl -sS "$BASE_URL/backend/media/multipart/parts?upload_id=$UPLOAD_ID" \
  -H "Authorization: $TOKEN" | python3 -m json.tool | head -20

echo "==> 4) presign + PUT"
PART_NUMS=$(python3 -c "print(','.join(str(i) for i in range(1, int('$PART_COUNT')+1)))")
curl -sS -X POST "$BASE_URL/backend/media/multipart/presign" \
  -H "Authorization: $TOKEN" -H 'Content-Type: application/json' \
  -d "{\"upload_id\":\"$UPLOAD_ID\",\"part_numbers\":[${PART_NUMS}]}" \
  > "$TMPDIR/presign.json"

python3 - "$FILE" "$PART_SIZE" "$TMPDIR/presign.json" "$TMPDIR/parts.json" <<'PY'
import json, sys, subprocess, pathlib
file_path, part_size = sys.argv[1], int(sys.argv[2])
presign = json.load(open(sys.argv[3]))
assert presign.get("code") == 0, presign
data = pathlib.Path(file_path).read_bytes()
etags = []
for item in presign["data"]["list"]:
    n = item["part_number"]
    start = (n - 1) * part_size
    chunk = data[start:start + part_size]
    out = subprocess.check_output([
        "curl", "-sS", "-D", "-", "-o", "/dev/null", "-X", "PUT",
        "-H", f"Content-Length: {len(chunk)}",
        "--data-binary", "@-",
        item["url"],
    ], input=chunk)
    headers = out.decode("utf-8", "ignore")
    etag = ""
    for line in headers.splitlines():
        if line.lower().startswith("etag:"):
            etag = line.split(":", 1)[1].strip().strip('"')
            break
    if not etag:
        raise SystemExit(f"part {n} 未返回 ETag:\n{headers[:800]}")
    print(f"    part {n}: {len(chunk)} bytes etag={etag}")
    etags.append({"part_number": n, "etag": etag})
json.dump(etags, open(sys.argv[4], "w"))
PY

echo "==> 5) complete"
PARTS_JSON=$(cat "$TMPDIR/parts.json")
curl -sS -X POST "$BASE_URL/backend/media/multipart/complete" \
  -H "Authorization: $TOKEN" -H 'Content-Type: application/json' \
  -d "{\"upload_id\":\"$UPLOAD_ID\",\"parts\":$PARTS_JSON}" \
  > "$TMPDIR/complete.json"
python3 - "$TMPDIR/complete.json" <<'PY'
import json,sys
d=json.load(open(sys.argv[1]))
assert d.get("code")==0, d
print("    id =", d["data"]["id"])
print("    url=", d["data"]["url"])
print("OK 上传成功")
PY
