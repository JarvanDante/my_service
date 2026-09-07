package kit

import "time"

// totpNow 可在单测里替换, 避免卡在 30 秒窗口边界。
var totpNow = time.Now
