#!/bin/sh
#
# Runs during git flow release start
#
# Positional arguments:
# $1    Version
#
# Return VERSION - When VERSION is returned empty gitflow 
#	will stop as the version is necessary
#
VERSION=$(cat .version)

# Implement your script here.
TAGS=`git tag v${VERSION}* -l|wc -l`
if [ "$TAGS" != 0 ]; then
	LASTTAG=$(git describe --tags $(git rev-list --tags --max-count=1) |tr -d v)
	MAJOR=`echo ${LASTTAG} | sed "s/^\([0-9]*\).*/\1/"`
	MINOR=`echo ${LASTTAG} | sed "s/[0-9]*\.\([0-9]*\).*/\1/"`
	REVISION=`git tag --list v$MAJOR.$MINOR.* |wc -l | tr -d ' '`
	VERSION=$MAJOR.$MINOR.$REVISION
fi
# Return the VERSION
echo ${VERSION} > .version
echo New version is v${VERSION}
exit 0
