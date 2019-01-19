

package console

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/YunWisdom/BookLog/service"
	"github.com/YunWisdom/BookLog/util"
	"github.com/gin-gonic/gin"
)

// ExportMarkdownAction exports articles as markdown zip file.
func ExportMarkdownAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	if nil == session {
		result.Code = -1
		result.Msg = "please login before export"

		return
	}

	tempDir := os.TempDir()
	logger.Trace("temp dir path is [" + tempDir + "]")
	zipFilePath := filepath.Join(tempDir, session.UName+"-export-md.zip")
	zipFile, err := util.Zip.Create(zipFilePath)
	if nil != err {
		logger.Errorf("create zip file [" + zipFilePath + "] failed: " + err.Error())
		result.Code = -1
		result.Msg = "create zip file failed"

		return
	}

	c.Header("Content-Disposition", "attachment; filename="+session.UName+"-export-md.zip")
	c.Header("Content-Type", "application/zip")

	mdFiles := service.Export.ExportMarkdowns(session.BID)
	if 1 > len(mdFiles) {
		zipFile.Close()
		file, err := os.Open(zipFilePath)
		if nil != err {
			logger.Errorf("open zip file [" + zipFilePath + " failed: " + err.Error())
			result.Code = -1
			result.Msg = "open zip file failed"

			return
		}
		defer file.Close()

		io.Copy(c.Writer, file)

		return
	}

	zipPath := filepath.Join(tempDir, session.UName+"-export-md")
	if err = os.RemoveAll(zipPath); nil != err {
		logger.Errorf("remove temp dir [" + zipPath + "] failed: " + err.Error())
		result.Code = -1
		result.Msg = "remove temp dir failed"

		return
	}
	if err = os.Mkdir(zipPath, 0755); nil != err {
		logger.Errorf("make temp dir [" + zipPath + "] failed: " + err.Error())
		result.Code = -1
		result.Msg = "make temp dir failed"

		return
	}
	for _, mdFile := range mdFiles {
		filename := filepath.Join(zipPath, mdFile.Name+".md")
		if err := ioutil.WriteFile(filename, []byte(mdFile.Content), 0644); nil != err {
			logger.Errorf("write file [" + filename + "] failed: " + err.Error())
		}
	}

	zipFile.AddDirectory(session.UName+"-export-md", zipPath)
	if err := zipFile.Close(); nil != err {
		logger.Errorf("zip failed: " + err.Error())
		result.Code = -1
		result.Msg = "zip failed"

		return
	}
	file, err := os.Open(zipFilePath)
	if nil != err {
		logger.Errorf("open zip file [" + zipFilePath + " failed: " + err.Error())
		result.Code = -1
		result.Msg = "open zip file failed"

		return
	}
	defer file.Close()

	io.Copy(c.Writer, file)
}
