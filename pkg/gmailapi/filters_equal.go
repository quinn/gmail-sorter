package gmailapi

import (
	"reflect"

	"google.golang.org/api/gmail/v1"
)

func filtersEqual(left, right *gmail.Filter) (match bool) {
	if left == nil || right == nil {
		return false
	}

	/*
		compare actions
	*/
	if left.Action.AddLabelIds == nil {
		left.Action.AddLabelIds = []string{}
	}

	if left.Action.RemoveLabelIds == nil {
		left.Action.RemoveLabelIds = []string{}
	}

	if right.Action.AddLabelIds == nil {
		right.Action.AddLabelIds = []string{}
	}

	if right.Action.RemoveLabelIds == nil {
		right.Action.RemoveLabelIds = []string{}
	}

	if !reflect.DeepEqual(left.Action.AddLabelIds, right.Action.AddLabelIds) {
		return false
	}

	if !reflect.DeepEqual(left.Action.RemoveLabelIds, right.Action.RemoveLabelIds) {
		return false
	}

	if left.Action.Forward != right.Action.Forward {
		return false
	}

	/*
		compare criteria
	*/
	if left.Criteria.ExcludeChats != right.Criteria.ExcludeChats {
		return false
	}

	if left.Criteria.From != right.Criteria.From {
		return false
	}

	if left.Criteria.HasAttachment != right.Criteria.HasAttachment {
		return false
	}

	if left.Criteria.NegatedQuery != right.Criteria.NegatedQuery {
		return false
	}

	if left.Criteria.Query != right.Criteria.Query {
		return false
	}

	if left.Criteria.Size != right.Criteria.Size {
		return false
	}

	if left.Criteria.SizeComparison != right.Criteria.SizeComparison {
		return false
	}

	if left.Criteria.Subject != right.Criteria.Subject {
		return false
	}

	if left.Criteria.To != right.Criteria.To {
		return false
	}

	return true
}
