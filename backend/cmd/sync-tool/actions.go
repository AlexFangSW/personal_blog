package main

import (
	"blog/entities"
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
)

func syncAll(ctx context.Context, baseURL, sourcePath string, batchSize int) error {
	slog.Info("syncAll")

	loginDone := make(chan bool, 1)
	processDone := make(chan bool, 1)
	processErr := make(chan error, 1)

	go func() {
		// login
		jwt, err := login(ctx, loginDone, baseURL)
		fmt.Printf("\n")
		if err != nil {
			processErr <- fmt.Errorf("syncAll: login failed: %w", err)
			return
		}
		slog.Debug("got jwt", "token", jwt)

		syncHelper := NewSyncHelper(baseURL, jwt, batchSize)

		// get data from server
		tags, err := syncHelper.GetAllTags()
		if err != nil {
			processErr <- fmt.Errorf("syncAll: failed to get tags from server: %w", err)
			return
		}

		topics, err := syncHelper.GetAllTopics()
		if err != nil {
			processErr <- fmt.Errorf("syncAll: failed to get topics from server: %w", err)
			return
		}

		blogs, err := syncHelper.GetAllBlogs()
		if err != nil {
			processErr <- fmt.Errorf("syncAll: failed to get blogs from server: %w", err)
			return
		}

		// load meta file
		metafile, err := loadMetaFile(filepath.Join(sourcePath, "meta.yaml"))
		if err != nil {
			processErr <- fmt.Errorf("syncAll: load meta file failed: %w", err)
			return
		}

		// load blogs
		idMap, err := loadIDMap(filepath.Join(sourcePath, "ids.json"))
		if err != nil {
			processErr <- fmt.Errorf("syncAll: load id map failed: %w", err)
			return
		}

		localblogs, err := loadBlogs(filepath.Join(sourcePath, "blogs"), idMap)
		if err != nil {
			processErr <- fmt.Errorf("syncAll: load blogs failed: %w", err)
			return
		}

		// seperate into groups (CRUD + noop)
		groupedTags, err := groupTags(metafile.Tags, tags)
		if err != nil {
			processErr <- fmt.Errorf("syncAll: group tags failed: %w", err)
			return
		}

		groupedTopics, err := groupTopics(metafile.Topics, topics)
		if err != nil {
			processErr <- fmt.Errorf("syncAll: group topics failed: %w", err)
			return
		}

		groupedBlogs, err := groupBlogs(localblogs, blogs)
		if err != nil {
			processErr <- fmt.Errorf("syncAll: group blogs failed: %w", err)
			return
		}

		// sync
		// create tags and topics, also fills in their ids for later use
		newTopics, err := syncHelper.CreateTopics(groupedTopics.create)
		if err != nil {
			processErr <- fmt.Errorf("syncAll: create topics failed: %w", err)
			return
		}

		newTags, err := syncHelper.CreateTags(groupedTags.create)
		if err != nil {
			processErr <- fmt.Errorf("syncAll: create tags failed: %w", err)
			return
		}

		// update tags and topics
		updatedTopics, err := syncHelper.UpdateTopics(groupedTopics.update)
		if err != nil {
			processErr <- fmt.Errorf("syncAll: update topics failed: %w", err)
			return
		}
		updatedTags, err := syncHelper.UpdateTags(groupedTags.update)
		if err != nil {
			processErr <- fmt.Errorf("syncAll: update tags failed: %w", err)
			return
		}

		// blogs
		// prepare blogs for CRUD operations
		existingTopics := slices.Concat[[]entities.Topic](newTopics, updatedTopics, groupedTopics.noop)
		existingTags := slices.Concat[[]entities.Tag](newTags, updatedTags, groupedTags.noop)
		blogTransformHelper := NewBlogTransformHelper(existingTags, existingTopics, sourcePath)
		transformedBlogs, err := blogTransformHelper.Transform(groupedBlogs)
		if err != nil {
			processErr <- fmt.Errorf("syncAll: transform blogs failed: %w", err)
			return
		}

		if err := syncHelper.CreateBlogs(transformedBlogs.create); err != nil {
			processErr <- fmt.Errorf("syncAll: create blogs failed: %w", err)
			return
		}
		if err := syncHelper.UpdateBlogs(transformedBlogs.update); err != nil {
			processErr <- fmt.Errorf("syncAll: create blogs failed: %w", err)
			return
		}
		if err := syncHelper.DeleteBlogs(transformedBlogs.delete); err != nil {
			processErr <- fmt.Errorf("syncAll: create blogs failed: %w", err)
			return
		}

		// delete tags and topics
		if err := syncHelper.DeleteTopics(groupedTopics.delete); err != nil {
			processErr <- fmt.Errorf("syncAll: delete topics failed: %w", err)
			return
		}
		if err := syncHelper.DeleteTags(groupedTags.delete); err != nil {
			processErr <- fmt.Errorf("syncAll: delete tags failed: %w", err)
			return
		}

		processDone <- true
	}()

	select {
	case <-ctx.Done():
		slog.Warn("got done")
		<-loginDone
		return nil
	case <-processDone:
		return nil
	case err := <-processErr:
		return err
	}
}
