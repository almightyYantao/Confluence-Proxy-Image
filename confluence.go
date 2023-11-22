package main

import (
	database "confluence-proxy-attachment/config"
	"fmt"
)

func confluencePath(pageId string, spaceId string, contentId string) string {
	attachmentPath := "ver003"

	for i := len(spaceId); i >= 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		substr := spaceId[start:i]
		result := modulo(substr, 250)
		attachmentPath = attachmentPath + "/" + fmt.Sprintf("%d", result)
		if start == len(spaceId)-6 {
			attachmentPath = attachmentPath + "/" + spaceId
			break
		}
	}

	for i := len(pageId); i >= 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		substr := pageId[start:i]
		result := modulo(substr, 250)
		attachmentPath = attachmentPath + "/" + fmt.Sprintf("%d", result)
		if start == len(pageId)-6 {
			attachmentPath = attachmentPath + "/" + pageId
			break
		}
	}
	return attachmentPath + "/" + contentId + "/1"
}

func modulo(substr string, module int) int {
	num := 0
	for _, ch := range substr {
		num = num*10 + int(ch-'0')
		num %= module
	}
	return num
}

func query(values ...interface{}) ConfluenceContent {
	if len(values) < 2 {
		return ConfluenceContent{}
	}
	var runSql string
	if values[0] == "title" && values[1] != nil && values[2] != nil {
		runSql = fmt.Sprintf("SELECT `CONTENTID`,`TITLE`,`PAGEID`,`SPACEID` FROM CONTENT WHERE `CONTENTTYPE` = 'ATTACHMENT' AND `TITLE` = '%s' AND `PAGEID` = '%s'", values[1], values[2])
	}
	if values[0] == "contentId" && values[1] != nil && values[2] != nil {
		runSql = fmt.Sprintf("SELECT `CONTENTID`,`TITLE`,`PAGEID`,`SPACEID` FROM CONTENT WHERE `CONTENTTYPE` = 'ATTACHMENT' AND `CONTENTID` = '%s' AND `PAGEID` = '%s'", values[1], values[2])
	}
	if len(runSql) > 0 {
		rows, err := database.DB.Query(runSql)
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()

		// Iterate through the results
		for rows.Next() {
			var confluenceContent ConfluenceContent
			if err := rows.Scan(&confluenceContent.CONTENTID, &confluenceContent.TITLE, &confluenceContent.PAGEID, &confluenceContent.SPACEID); err != nil {
				panic(err.Error())
			}
			return confluenceContent
		}

		if err = rows.Err(); err != nil {
			panic(err.Error())
		}
	}
	return ConfluenceContent{}
}
