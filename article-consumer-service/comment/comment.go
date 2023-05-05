package comment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/pubsub"
)

func DeleteByArticleID(ctx context.Context, msg pubsub.SubscribeDater, client lib.CustomArticleClient[lib.ArticleComment]) error {
	event := pubsub.DefaultEvent{}
	err := json.Unmarshal(msg.Data(), &event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	if event.ActionType != pubsub.ActionTypeDelete {
		// skip if not delete event
		return nil
	}

	err = client.DeleteByArticleID(ctx, lib.ArticleComment{
		ArticleID: event.ResourceID,
	})
	if err != nil {
		return fmt.Errorf("delete by article id returned an error: %w", err)
	}
	msg.Ack()
	return nil
}
