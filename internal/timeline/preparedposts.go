package timeline

import (
	"container/list"
	"errors"

	apimodel "github.com/superseriousbusiness/gotosocial/internal/api/model"
)

type preparedPosts struct {
	data *list.List
}

type preparedPostsEntry struct {
	statusID string
	prepared *apimodel.Status
}

func (p *preparedPosts) insertPrepared(i *preparedPostsEntry) error {
	if p.data == nil {
		p.data = &list.List{}
	}

	// if we have no entries yet, this is both the newest and oldest entry, so just put it in the front
	if p.data.Len() == 0 {
		p.data.PushFront(i)
		return nil
	}

	var insertMark *list.Element
	// We need to iterate through the index to make sure we put this post in the appropriate place according to when it was created.
	// We also need to make sure we're not inserting a duplicate post -- this can happen sometimes and it's not nice UX (*shudder*).
	for e := p.data.Front(); e != nil; e = e.Next() {
		entry, ok := e.Value.(*preparedPostsEntry)
		if !ok {
			return errors.New("index: could not parse e as a preparedPostsEntry")
		}

		// if the post to index is newer than e, insert it before e in the list
		if insertMark == nil {
			if i.statusID > entry.statusID {
				insertMark = e
			}
		}

		// make sure we don't insert a duplicate
		if entry.statusID == i.statusID {
			return nil
		}
	}

	if insertMark != nil {
		p.data.InsertBefore(i, insertMark)
		return nil
	}

	// if we reach this point it's the oldest post we've seen so put it at the back
	p.data.PushBack(i)
	return nil
}