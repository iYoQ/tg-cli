package app

import "testing"

func TestChatItem(t *testing.T) {
	item := chatItem{title: "general", id: 42, haveTopics: true}

	if got := item.Title(); got != "general" {
		t.Errorf("Title() = %q, want %q", got, "general")
	}
	if got := item.Description(); got != "ID:42" {
		t.Errorf("Description() = %q, want %q", got, "ID:42")
	}
	if got := item.FilterValue(); got != "general" {
		t.Errorf("FilterValue() = %q, want %q", got, "general")
	}
}

func TestTopicItem(t *testing.T) {
	item := topicItem{chatId: 7, threadId: 11, title: "random"}

	if got := item.Title(); got != "random" {
		t.Errorf("Title() = %q, want %q", got, "random")
	}
	if got := item.Description(); got != "ID:7" {
		t.Errorf("Description() = %q, want %q", got, "ID:7")
	}
	if got := item.FilterValue(); got != "random" {
		t.Errorf("FilterValue() = %q, want %q", got, "random")
	}
}
