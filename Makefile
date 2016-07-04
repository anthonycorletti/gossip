PROG_NAME=gossip
OUT_DIR=bin
DEFAULT_INSTALL=/usr/local/bin

LESS_FILES=$(wildcard client/css/*.less)
CSS_FILES=$(LESS_FILES:.less=.css)
CLIENT_JS_FILES=$(wildcard client/js/*.js)
HTML_FILES=client/views/*

VAGRANT_KEY=~/.vagrant.d/insecure_private_key
VAGRANT_INVENTORY=./.vagrant/provisioners/ansible/inventory/vagrant_ansible_inventory

ANS_ENV_LOCAL=ANSIBLE_HOST_KEY_CHECKING=False
ANS_OPTS_LOCAL=--private-key=$(VAGRANT_KEY) --user=vagrant -i $(VAGRANT_INVENTORY)
LOCAL_ANSIBLE=$(ANS_ENV_LOCAL) ansible-playbook $(ANS_OPTS_LOCAL)

LESS=./node_modules/less/bin/lessc
MINIFY=node ./node_modules/minifier/index.js


default: sane-output deps fmt binary client

sane-output:
	mkdir -p $(OUT_DIR)
	mkdir -p $(OUT_DIR)/views
	mkdir -p $(OUT_DIR)/js
	mkdir -p $(OUT_DIR)/css

binary:
	@go build -o $(OUT_DIR)/$(PROG_NAME) ./src/
	@echo "***** Binary Built"

run: default
	@$(OUT_DIR)/gossip -d $(OUT_DIR) -l ./logs

deps:
	@depman
	@npm install
	@echo "***** Dependencies Met"

clean:
	@rm -rf logs $(OUT_DIR) node_modules
	@go clean ./src/

todo:
	@grep -nri "todo"


install: default directories install-binary install-client

directories:
	@mkdir -p /srv/$(PROG_NAME)

stop-service:
	@-service gossip stop

install-binary: stop-service
	@if test "$(PREFIX)" = "" ; then \
		cp $(OUT_DIR)/$(PROG_NAME) $(DEFAULT_INSTALL)/$(PROG_NAME) ; \
	else \
		cp $(OUT_DIR)/$(PROG_NAME) $(PREFIX)/$(PROG_NAME); \
	fi

install-client:
	@cp -r $(OUT_DIR)/* /srv/$(PROG_NAME)/


fmt:
	@go fmt ./src/


client: stylesheet javascript static-js html
	@cp client/config.json $(OUT_DIR)/
	@echo "***** Client Finished"

stylesheet: $(CSS_FILES)
	@echo "***** LESS Compiled"
	@cp client/css/*.css $(OUT_DIR)/css/

%.css: %.less
	@$(LESS) $< > $(OUT_DIR)/css/$(notdir $@)

clean-js:
	@rm -f $(OUT_DIR)/js/app.js

javascript: clean-js $(OUT_DIR)/js/app.js
	@echo "***** Javascript Merged and Minified"

# $(OUT_DIR)/js/app.js: $(CLIENT_JS_FILES)

$(OUT_DIR)/js/app.js:
	@cat $(CLIENT_JS_FILES) >> $(OUT_DIR)/js/app.js

$(OUT_DIR)/js/app.min.js: $(OUT_DIR)/js/app.js
	@$(MINIFY) --output $(OUT_DIR)/js/app.min.js $(OUT_DIR)/js/app.js

html:
	@cp $(HTML_FILES) $(OUT_DIR)/views/

static-js:
	@cp client/deps/*.js $(OUT_DIR)/js/


vm:
	@vagrant up

local: vm
	$(LOCAL_ANSIBLE) ansible/all.yml

local-db: vm
	$(LOCAL_ANSIBLE) ansible/db.yml

local-api: vm
	@$(LOCAL_ANSIBLE) ansible/api.yml
