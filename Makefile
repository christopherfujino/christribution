GIT_REMOTE = https://github.com/lfs-book/lfs.git
# Aug 31, 2025
TAG = r12.4

third_party/LFS-BOOK.html:
	# Make a single file html doc
	cd third_party/lfs && make BASEDIR=.. nochunks

# clone LFS
third_party/lfs/README: third_party/.stamp
	git clone $(GIT_REMOTE) -b $(TAG) third_party/lfs

third_party/.stamp:
	mkdir -p third_party
	touch third_party/.stamp

.PHONY: lfs-build-deps
lfs-build-deps:
	apt-get install libxml2 libxml2-utils libxslt1.1 docbook5-xml docbook-xsl libtidy58 xsltproc tidy
