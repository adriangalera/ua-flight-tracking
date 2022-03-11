PACKAGE_NAME=flight-tracking
BUILD_DIR=build/
BUILD_FILE=$(BUILD_DIR)$(PACKAGE_NAME)-${platform}-${version}

make: build

build:
	go build -o $(BUILD_FILE)

clean:
	rm -f $(BUILD_FILE)

run:
	$(BUILD_FILE)