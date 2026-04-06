
.PHONY: all

all: package package2

package:
	env COPYFILE_DISABLE=1 tar -czvf ../agent_src.tar.gz Dockerfile requirements.txt *.py

package2:
	env COPYFILE_DISABLE=1 tar -czvf ../agent_src2.tar.gz Dockerfile requirements.txt *.py
