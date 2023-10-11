#pragma once
#include <headers.h>

void Log(std::string msg,
         const std::source_location location = std::source_location::current());

template <typename... Args>
void TraceLog(const std::source_location location, Args &&...args) {
  std::ostringstream stream;
  stream << location.file_name() << ":" << location.line() << " : ";
  (stream << ... << std::forward<Args>(args)) << '\n';

  std::cout << stream.str();
}