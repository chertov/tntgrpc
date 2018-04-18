#include "errors.hpp"

#include <iostream>

GRPCException::GRPCException(bool noerr) : noerr_(noerr) {}
GRPCException::GRPCException(char const* fmt, ...) {
    char text[1000];
    va_list ap;
    va_start(ap, fmt);
    vsnprintf(text, sizeof text, fmt, ap);
    va_end(ap);
    ru.append(text);
    en.append(text);
}

std::string GRPCException::log() const {
    switch (lang()) {
        case EN: { return en; }
        case RU: { return ru; }
    }
    return "";
}
std::string GRPCException::str(int level) const {
    std::stringstream ss;
    ss << location << "    Code: " << code() << "    Type: " << type() << "    " << log();
    if(lastError) ss << std::endl << std::string(level*4, ' ') << lastError->str(level+1);
    return ss.str();
}

GRPCError GRPCException::run(const std::string location) {
    this->location = location;
    return shared_from_this();
}
GRPCError GRPCException::run(const std::string location, const GRPCError &lastError) {
    this->lastError = lastError;
    return run(location);
}