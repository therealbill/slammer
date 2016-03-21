VERSION = $(shell cat .version)
GHT = $(GITHUB_TOKEN)
LAST_TAG = $(shell  git describe --abbrev=0 --tags | tr -d v)

all: 
	@echo Building firehose
	go build -o firehose/firehose -ldflags "-X main.VERSION=${VERSION}" firehose/main.go 
	@echo Building firetruck
	go build -o firetruck/firetruck -ldflags "-X main.VERSION=${VERSION}" firetruck/main.go 

release: 
ifeq "$(VERSION)" "$(LAST_TAG)"
	@echo "No Version bump found, not doing anything. You'll need to run 'make major', 'make minor', or 'make revision' depending on how much you want to bump the version"
else 
	@echo Building release v$(VERSION)
	@echo Creating Git tag
	@git tag -a v$(shell cat .version) -m "Version $(shell cat .version)"
	@echo Pushing git tag
	@git push && git push --tags
endif

revision:
	@echo Bumping revision. Last version was $(VERSION)
	@./scripts/bump-revision.sh
	@echo adding .version to git commit
	@git add .version
	@echo Now ready to run 'make release' to push the tag to github which will trigger travis-ci.org to build the release if all is well

minor:
	@echo Bumping minor. Last version was $(VERSION)
	@./scripts/bump-minor.sh
	@echo adding .version to git commit
	@git add .version
	@echo Now ready to run 'make release' to push the tag to github which will trigger travis-ci.org to build the release if all is well

major:
	@echo Bumping major. Last version was $(VERSION)
	@./scripts/bump-major.sh
	@echo adding .version to git commit
	@git add .version
	@echo Now ready to run 'make release' to push the tag to github which will trigger travis-ci.org to build the release if all is well

