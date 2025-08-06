// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package localos

import (
	"fmt"
	"os"
)

func GetLocalOSHost() string {
	return fmt.Sprintf("%s:%s", os.Getenv("COZE_LOOP_OSS_DOMAIN"), os.Getenv("COZE_LOOP_OSS_PORT"))
}
