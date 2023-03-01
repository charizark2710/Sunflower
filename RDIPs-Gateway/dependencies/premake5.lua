project "amqcpp"
    kind "StaticLib"
    language "C++"
    targetdir("./bin/%{prj.name}")
    cppdialect "C++17"
    objdir("./bin-int/%{prj.name}")
    files {"./AMQP-CPP/include/**.h", "./AMQP-CPP/src/**.cpp", "./AMQP-CPP/src/**.h"}

    includedirs {
        "./AMQP-CPP/include",
        "./AMQP-CPP/src"
    }

    filter "configurations:Debug"
        defines {"DEBUG"}
        symbols "On"

    filter "configurations:Release"
        defines {"NDEBUG"}
        optimize "On"
