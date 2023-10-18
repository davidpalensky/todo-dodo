package action

import (
	"fmt"
	"log"
	"todo-dodo/db"
)

// Inserts a bunch of tags into the database, inserts into task_tag_links aswell, if task_id_link is not nil
func TagCreateBatch(tags []TagFetcher, task_id_link *uint64) error {
	if task_id_link != nil {
		return tagCreateBatchLinked(tags, *task_id_link)
	}
	for _, tag := range tags {
		_, err := db.DB.Exec("INSERT INTO tags (user_id, title, color) VALUES (?, ?, ?) ON CONFLICT (user_id, title, color) DO NOTHING;", tag.User_id, tag.Title, tag.Color)
		if err != nil {
			//db.DB.Exec("ROLLBACK;")
			//log.Printf("Error: Could not insert data into db: %s", err)
			return &ActionError{Kind: "database", Msg: "TaskCreate: Could not enter tags into db: " + err.Error()}
		}
	}
	return nil
}

func tagCreateBatchLinked(tags []TagFetcher, task_id uint64) error {
	for _, tag := range tags {
		// Yes, this is awkward af
		var tag_id_retriever []uint64
		err := db.DB.Select(&tag_id_retriever, "SELECT tag_id FROM tags WHERE user_id = ? AND title = ? AND color = ?;", tag.User_id, tag.Title, tag.Color)
		if err != nil {
			//log.Printf("Error: Could not insert data into db: %s", err)
			log.Fatalf("This query should not fail: %s", err.Error())
		}
		// If tag is not in there, add that shit
		if len(tag_id_retriever) == 0 {
			res, err := db.DB.Exec("INSERT INTO tags (user_id, title, color) VALUES (?, ?, ?) ON CONFLICT (user_id, title, color) DO NOTHING;", tag.User_id, tag.Title, tag.Color)
			tag_id_int64, _ := res.LastInsertId()
			tag_id := uint64(tag_id_int64)
			if err != nil {
				//log.Printf("Error: Could not insert data into db: %s", err)
				return &ActionError{Kind: "database", Msg: "TaskCreate: Could not enter tags into db: " + err.Error()}
			}
			// Insert link
			_, err1 := db.DB.Exec("INSERT INTO task_tag_links (task_id, tag_id) VALUES (?, ?);", task_id, tag_id)
			if err1 != nil {
				log.Printf("Error: Could not insert data into db: task_id = %d, tag_id = %d", task_id, tag_id)
				return &ActionError{Kind: "database", Msg: "TaskCreate: Could not enter task tag links into db: " + err1.Error()}
			}
		} else { // Insert link
			tag_id := tag_id_retriever[0]
			_, err := db.DB.Exec("INSERT INTO task_tag_links (task_id, tag_id) VALUES (?, ?);", task_id, tag_id)
			if err != nil {
				log.Printf("Error: Could not insert data into db: task_id = %d, tag_id = %d", task_id, tag_id)
				return &ActionError{Kind: "database", Msg: "TaskCreate: Could not enter task tag links into db: " + err.Error()}
			}
		}
	}
	return nil
}

// Required data for creating a tag
type TagFetcher struct {
	User_id uint64 `json:"user_id"`
	Title   string `json:"title"`
	Color   string `json:"color"`
}

// Required data for creating a tag
type TagModel struct {
	Tag_id  uint64 `json:"tag_id"`
	User_id uint64 `json:"user_id"`
	Title   string `json:"title"`
	Color   string `json:"color"`
}

// Required data for creating a tag
type TagFetchAllArgs struct {
	User_id uint `json:"user_id"`
}

// Model of the tag data in the DB
type TagFetchAllDBReturn struct {
	Tag_id  uint64 `json:"tag_id"`
	User_id uint64 `json:"user_id"`
	Title   string `json:"title"`
	Color   string `json:"color"`
}

// Fetch all tags
func TagFetchAll(args *TagFetchAllArgs) ([]TagFetchAllDBReturn, error) {
	var result []TagFetchAllDBReturn
	err := db.DB.Select(&result, "SELECT * FROM tags WHERE user_id = ?;", fmt.Sprintf("%d", args.User_id))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete a batch of tags
func TagDeleteBatch(tag_ids []uint64) error {
	for _, tag_id := range tag_ids {
		_, err1 := db.DB.Exec("DELETE FROM task_tag_links WHERE tag_id = ?;", tag_id)
		if err1 != nil {
			return &ActionError{Kind: "database", Msg: "Unable to delete tag with task_id " + fmt.Sprintf("%d", tag_id)}
		}
		_, err := db.DB.Exec("DELETE FROM tags WHERE tag_id = ?;", tag_id)
		if err != nil {
			return &ActionError{Kind: "database", Msg: "Unable to delete tag with task_id " + fmt.Sprintf("%d", tag_id)}
		}
	}
	return nil
}
