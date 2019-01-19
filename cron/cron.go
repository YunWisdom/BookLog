

// Package cron includes all cron tasks.
package cron

import (
	"os"

	"github.com/YunWisdom/BookLog/log"
)

// Logger
var logger = log.NewLogger(os.Stdout)

// Start starts all cron tasks.
func Start() {
	refreshRecommendArticlesPeriodically()
	pushArticlesPeriodically()
	pushCommentsPeriodically()
}
