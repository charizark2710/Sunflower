
#include <string>
#include <utils/Log.h>

void Log(std::string msg, const std::source_location location) {
  std::cout << location.file_name() << ":" << location.line() << " : " << msg << "\n";
}
