package logic

import (
	"fmt"
	"log"
	"todo-dodo/db"
	"todo-dodo/util"
)

// Inserts a bunch of tags into the database, inserts into task_tag_links aswell, if task_id_link is not nil
func TagCreateBatch(tags []TagCreator, task_id_link *uint64) error {
	if task_id_link != nil {
		return tagCreateBatchLinked(tags, *task_id_link)
	}
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		_, err := tx.Exec("INSERT INTO tags (user_id, title, color) VALUES (?, ?, ?) ON CONFLICT (user_id, title, color) DO NOTHING;", tag.User_id, tag.Title, tag.Color)
		if err != nil {
			//db.DB.Exec("ROLLBACK;")
			//log.Printf("Error: Could not insert data into db: %s", err)
			tx.Rollback()
			return &LogicError{Kind: "database", Msg: "TaskCreate: Could not enter tags into db: " + err.Error()}
		}
	}
	tx.Commit()
	return nil
}

// Yes this function is very messy and does too many queries, i am not a sqlite magician however.
func tagCreateBatchLinked(tags []TagCreator, task_id uint64) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		// Verify color
		if !util.ValidateHexcode(tag.Color) {
			tx.Rollback()
			return &LogicError{Kind: "invalid data", Msg: "Color code `" + tag.Color + "` is not a valid hex color code."}
		}

		// Yes, this is awkward af
		var tag_id_retriever []uint64
		err := tx.Select(&tag_id_retriever, "SELECT tag_id FROM tags WHERE user_id = ? AND title = ? AND color = ?;", tag.User_id, tag.Title, tag.Color)
		if err != nil {
			//log.Printf("Error: Could not insert data into db: %s", err)
			log.Fatalf("This query should not fail: %s", err.Error())
		}
		// If tag is not in there, add that shit
		if len(tag_id_retriever) == 0 {
			res, err := tx.Exec("INSERT INTO tags (user_id, title, color) VALUES (?, ?, ?);", tag.User_id, tag.Title, tag.Color)
			tag_id_int64, _ := res.LastInsertId()
			tag_id := uint64(tag_id_int64)
			if err != nil {
				//log.Printf("Error: Could not insert data into db: %s", err)
				tx.Rollback()
				return err
			}
			// Insert link
			_, err1 := tx.Exec("INSERT INTO task_tag_links (task_id, tag_id) VALUES (?, ?);", task_id, tag_id)
			if err1 != nil {
				//log.Printf("Error: Could not insert data into db: task_id = %d, tag_id = %d", task_id, tag_id)
				tx.Rollback()
				return err1
			}
		} else { // Insert link
			tag_id := tag_id_retriever[0]
			_, err := tx.Exec("INSERT INTO task_tag_links (task_id, tag_id) VALUES (?, ?);", task_id, tag_id)
			if err != nil {
				//log.Printf("Error: Could not insert data into db: task_id = %d, tag_id = %d", task_id, tag_id)
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

// Required data for creating a tag
type TagCreator struct {
	User_id uint64 `json:"user_id"`
	Title   string `json:"title"`
	Color   string `json:"color"`
}

// Model of the tag data in the DB
type TagModel struct {
	Tag_id  uint64 `json:"tag_id"`
	User_id uint64 `json:"user_id"`
	Title   string `json:"title"`
	Color   string `json:"color"`
}

// Required data for creating a tag
type TagFetchAllArgs struct {
	User_id uint64 `json:"user_id"`
}

// Fetch all tags
func TagFetchAll(a *TagFetchAllArgs) ([]TagModel, error) {
	var result []TagModel
	err := db.DB.Select(&result, "SELECT * FROM tags WHERE user_id = ?;", fmt.Sprintf("%d", a.User_id))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete a batch of tags
func TagDeleteBatch(tag_ids []uint64) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	for _, tag_id := range tag_ids {
		_, err1 := tx.Exec("DELETE FROM task_tag_links WHERE tag_id = ?;", tag_id)
		if err1 != nil {
			return &LogicError{Kind: "database", Msg: "Unable to delete tag with task_id " + fmt.Sprintf("%d", tag_id)}
		}
		_, err := tx.Exec("DELETE FROM tags WHERE tag_id = ?;", tag_id)
		if err != nil {
			return &LogicError{Kind: "database", Msg: "Unable to delete tag with task_id " + fmt.Sprintf("%d", tag_id)}
		}
	}
	tx.Commit()
	return nil
}
