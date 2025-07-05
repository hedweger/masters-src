cmake_minimum_required(VERSION 3.0)
project(ied-client C)

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
set(CMAKE_EXPORT_COMPILE_COMMANDS 1)
set(C_STANDARD c99)

include_directories(
	.
	/usr/local/include/libiec61850
)
link_directories(/usr/local/lib)

add_executable(ied-client client.c)
target_link_libraries(ied-client iec61850)
set_target_properties(ied-client PROPERTIES
	INSTALL_RPATH "/usr/local/lib"
	BUILD_WITH_INSTALL_RPATH TRUE
)
