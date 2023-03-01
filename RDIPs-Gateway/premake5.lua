workspace "RDIPs-Gateway"
    configurations {"Debug", "Release"}
    location "build"
    architecture "x64"

outputDir = "%{cfg.buildcfg}-%{cfg.system}-%{cfg.architecture}"

Includedirs = {}
Includedirs["amqcpp"] =  "./dependencies/AMQP-CPP/include"


include "./dependencies"

project "RDIPs-Gateway"
    kind "ConsoleApp"
    language "C++"
    cppdialect "C++17"
    targetdir("./bin")
    objdir("./bin-int")
    files {"./src/**.h", "./src/**.cpp", "./include/**,h"}
    
    includedirs {
        "%{Includedirs.amqcpp}",
        "./include",
        "./src"
    }

    links {
        "amqcpp",
        "pthread",
        "dl"
    }

    filter "configurations:Debug"
        defines {"DEBUG"}
        symbols "On"

    filter "configurations:Release"
        defines {"NDEBUG"}
        optimize "On"