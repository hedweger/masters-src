cmake_minimum_required(VERSION 3.5)
project(foo C)

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
set(CMAKE_EXPORT_COMPILE_COMMANDS 1)

include_directories(
	.
	/opt/homebrew/include/libiec61850
)
link_directories(/opt/homebrew/lib)

add_executable(foo main.c)

target_link_libraries(foo iec61850)
set_target_properties(foo PROPERTIES
    INSTALL_RPATH "/usr/local/lib"
    BUILD_WITH_INSTALL_RPATH TRUE
)

