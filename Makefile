tools = reference-gen

%:
	@ for tool in ${tools}; do \
	  make -C $${tool} $* ; \
	done
