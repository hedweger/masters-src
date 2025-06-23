cmake_minimum_required(VERSION 4.0)
project(iec61850 C)

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
set(CMAKE_EXPORT_COMPILE_COMMANDS 1)

include_directories(
	.
	/opt/homebrew/include/libiec61850
)
link_directories(/opt/homebrew/lib)

if(NOT DEFINED TARGET)
	add_executable(server server.c)
	add_executable(client client.c)
	foreach(tgt IN ITEMS server client)
		target_link_libraries(${tgt} iec61850)
		set_target_properties(${tgt} PROPERTIES
			INSTALL_RPATH "/usr/local/lib"
			BUILD_WITH_INSTALL_RPATH TRUE
		)
	endforeach()
else ()
if(TARGET STREQUAL "server")
	set(SOURCE_FILE server.c)
elseif(TARGET STREQUAL "client")
	set(SOURCE_FILE client.c)
endif()
add_executable(${TARGET} ${SOURCE_FILE})
target_link_libraries(${TARGET} iec61850)
set_target_properties(${TARGET} PROPERTIES
	INSTALL_RPATH "/usr/local/lib"
	BUILD_WITH_INSTALL_RPATH TRUE
)
endif()

