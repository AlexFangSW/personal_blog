package main

import (
	"fmt"
	"log/slog"
)

func syncAll(baseURL, metaFile, blogsDir string) error {
	// login
	jwt, err := login(baseURL)
	if err != nil {
		return fmt.Errorf("syncAll: login failed: %w", err)
	}
	slog.Info("got jwt", "token", jwt)

	syncHelper := NewSyncHelper(baseURL, jwt)

	// get data from server
	tags, err := syncHelper.GetAllTags()
	if err != nil {
		return fmt.Errorf("syncAll: failed to get tags from server: %w", err)
	}
	slog.Info("got tags", "tags", tags)

	topics, err := syncHelper.GetAllTopics()
	if err != nil {
		return fmt.Errorf("syncAll: failed to get topics from server: %w", err)
	}
	blogs, err := syncHelper.GetAllBlogs()
	if err != nil {
		return fmt.Errorf("syncAll: failed to get blogs from server: %w", err)
	}

	// load meta file
	metafile, err := loadMetaFile(metaFile)
	if err != nil {
		return fmt.Errorf("syncAll: load meta file failed: %w", err)
	}
	// load blogs
	localblogs, err := loadBlogs(blogsDir)
	if err != nil {
		return fmt.Errorf("syncAll: load blogs failed: %w", err)
	}

	// seperate into groups (CRUD + noop)
	groupedTags, err := groupTags(metafile.tags, tags)
	if err != nil {
		return fmt.Errorf("syncAll: group tags failed: %w", err)
	}
	groupedTopics, err := groupTopics(metafile.topics, topics)
	if err != nil {
		return fmt.Errorf("syncAll: group topics failed: %w", err)
	}
	groupedBlogs, err := groupBlogs(localblogs, blogs)
	if err != nil {
		return fmt.Errorf("syncAll: group blogs failed: %w", err)
	}

	// after this state blogs only lack content
	groupedInBlogs, err := transformBlogs(groupedTags, groupedTopics, groupedBlogs)
	if err != nil {
		return fmt.Errorf("syncAll: transform blogs failed: %w", err)
	}

	// sync

	// create tags and topics
	if err := syncHelper.CreateTopics(groupedTopics.create); err != nil {
		return fmt.Errorf("syncAll: create topics failed: %w", err)
	}
	if err := syncHelper.CreateTags(groupedTags.create); err != nil {
		return fmt.Errorf("syncAll: create tags failed: %w", err)
	}

	// update tags and topics
	if err := syncHelper.UpdateTopics(groupedTopics.update); err != nil {
		return fmt.Errorf("syncAll: update topics failed: %w", err)
	}
	if err := syncHelper.UpdateTags(groupedTags.update); err != nil {
		return fmt.Errorf("syncAll: update tags failed: %w", err)
	}

	// blogs
	if err := syncHelper.CreateBlogs(groupedInBlogs.create); err != nil {
		return fmt.Errorf("syncAll: create blogs failed: %w", err)
	}
	if err := syncHelper.UpdateBlogs(groupedInBlogs.update); err != nil {
		return fmt.Errorf("syncAll: create blogs failed: %w", err)
	}
	if err := syncHelper.DeleteBlogs(groupedInBlogs.delete); err != nil {
		return fmt.Errorf("syncAll: create blogs failed: %w", err)
	}

	// delete tags and topics
	if err := syncHelper.DeleteTopics(groupedTopics.delete); err != nil {
		return fmt.Errorf("syncAll: delete topics failed: %w", err)
	}
	if err := syncHelper.DeleteTags(groupedTags.delete); err != nil {
		return fmt.Errorf("syncAll: delete tags failed: %w", err)
	}

	return nil
}
