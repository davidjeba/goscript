# Advanced Capabilities Checklist

This checklist separates what now exists as a foundation in the repo from what still needs full product-grade implementation.

## 1. Vibe motion system

### 1.1 Added now

- [x] Transition primitives for spring, tween, and inertia
- [x] Motion values with subscriptions and velocity
- [x] Variants and gesture target mapping
- [x] Presence bookkeeping for exit flows
- [x] Layout delta calculation for transform-based motion
- [x] Scroll progress and in-view option models

### 1.2 Still missing

- [ ] Bind Vibe to real GoUIX / GoScript rendering
- [ ] Runtime animation execution on the client
- [ ] Shared element transitions between views
- [ ] Scroll observers and viewport tracking
- [ ] Drag physics and momentum
- [ ] Staggers, timelines, and orchestration
- [ ] Reduced-motion support

## 2. Internationalization

### 2.1 Added now

- [x] Locale bundle and catalog model
- [x] Message lookup with fallback locale support
- [x] Variable interpolation
- [x] RTL / LTR direction detection

### 2.2 Still missing

- [ ] ICU-style pluralization and rich message formatting
- [ ] Date, time, and number formatting helpers
- [ ] Route-level locale negotiation
- [ ] Resource file loading from disk
- [ ] Locale-aware input and validation
- [ ] Full RTL-aware component/layout behavior

## 3. Notion-style docs engine

### 3.1 Added now

- [x] Block-based document model
- [x] Nested page/content tree support
- [x] Inline database table schema model
- [x] File-browser tree model

### 3.2 Still missing

- [ ] Rich text spans and inline marks
- [ ] Collaborative block editing
- [ ] Slash commands and block transforms
- [ ] Nested toggles, synced blocks, embeds, mentions
- [ ] Comments, suggestions, and change tracking
- [ ] Database views like board, calendar, gallery, timeline

## 4. Inline editable tables

### 4.1 Added now

- [x] Editable column metadata
- [x] Row/cell updates with schema checks

### 4.2 Still missing

- [ ] Formula columns
- [ ] Rollups and relations
- [ ] Filter, sort, and grouping
- [ ] Pagination and virtualization
- [ ] Cell editors for date/select/file/relation types
- [ ] Live sync with backend records

## 5. File browser and workspace shell

### 5.1 Added now

- [x] File/directory tree model
- [x] Path lookup

### 5.2 Still missing

- [ ] Real filesystem adapter
- [ ] Permissions and policy hooks
- [ ] Upload, rename, move, trash, restore
- [ ] Previewers for image, text, markdown, and code
- [ ] Search, breadcrumbs, favorites, and recent files

## 6. Notion-class application platform

- [ ] Docs pages rendered entirely in GoScript
- [ ] Inline editing surface with block handles
- [ ] Database-backed pages and linked views
- [ ] Sidebar tree, breadcrumbs, tabs, and command palette
- [ ] Drag-and-drop between blocks, tables, and files
- [ ] Presence and multiplayer collaboration
- [ ] AI-assisted editing, summarization, and generation

## 7. Recommended build order

1. Bind `pkg/vibe` into GoUIX and GoScript components
2. Add `pkg/i18n` integration to routing, rendering, and UI direction
3. Build a workspace UI layer on top of `pkg/workspace`
4. Implement inline table views and file browser components
5. Add collaboration, syncing, and AI editor workflows
