workspace "RDIPs-Gateway"
    configurations {"Debug", "Release"}
    location "build"
    architecture "x64"

include "./dependencies/amqp.lua"
include "./dependencies/json.lua"

outputDir = "%{cfg.buildcfg}-%{cfg.system}-%{cfg.architecture}"

Includedirs = {}
Includedirs["amqcpp"] = "./dependencies/AMQP-CPP/include"
Includedirs["json"] = "./dependencies/json/include"

project "RDIPs-Gateway"
    kind "ConsoleApp"
    language "C++"
    cppdialect "C++20"
    targetdir("./bin")
    objdir("./bin-int")
    files {"./src/**.h", "./src/**.cpp", "./include/**.h"}
    
    includedirs {
        "%{Includedirs.amqcpp}",
        "./include",
        "./src"
    }

    links {
        "amqcpp",
        "json",
        "pthread",
        "dl"
    }

    filter "configurations:Debug"
        defines {"DEBUG"}
        symbols "On"

    filter "configurations:Release"
        defines {"NDEBUG"}
        optimize "On"