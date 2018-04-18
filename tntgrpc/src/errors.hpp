#pragma once

#include <sstream>
#include <vector>
#include <iostream>
#include <cstdarg>
#include "json.hpp"

enum ErrorLang { EN, RU };
inline ErrorLang lang() { return EN; }

#define SRC_POS std::string(__FILE__) + ":" + std::to_string(__LINE__)

class GRPCException;
typedef std::shared_ptr<GRPCException> GRPCError;
#define NOERR GRPCException::New(true)

#define grpcTry(code, grpcException) \
    { GRPCError err = [&](...) -> GRPCError { code; return NOERR; }(); if(!err->noerr()) return grpcException->run(SRC_POS, err); }

#define grpcIfErr(code, grpcException) \
    { GRPCError err = code; if(!err->noerr()) return grpcException->run(SRC_POS, err); }

#define grpcIfErrCheck(code) \
    { GRPCError err = code; if(!err->noerr()) return err->run(SRC_POS); }

class GRPCException : public std::enable_shared_from_this<GRPCException> {
protected:
    GRPCException(bool noerr = false);
    GRPCException(char const* fmt, ...);

public:
    static GRPCError New(bool noerr = false) { return std::make_shared<GRPCException>(GRPCException(noerr)); }
    static GRPCError New(char const* fmt, ...) {
        char text[1000];
        va_list ap;
        va_start(ap, fmt);
        vsnprintf(text, sizeof text, fmt, ap);
        va_end(ap);
        return std::make_shared<GRPCException>(GRPCException(text));
    }
    virtual int code() const { return 0; }
    virtual std::string type() const { return "GRPCError"; }
    virtual std::string log() const;
    virtual std::string str(int level = 1) const;
    bool noerr() const  { return noerr_; }
    GRPCError run(const std::string location = "");
    GRPCError run(const std::string location, const GRPCError &lastError);

    std::string json() const {
        nlohmann::json stack = make_stack();
        nlohmann::json current_err = stack[0];
        current_err["stack"] = stack;
        return current_err.dump(4);
    }
    virtual nlohmann::json to_json() const { return {{"errorCode",code()},{"log",log()},{"loc",location}}; }
    nlohmann::json make_stack() const {
        nlohmann::json stack = nlohmann::json::array();
        if(lastError && !lastError->noerr()) {
            nlohmann::json parent_stack = lastError->make_stack();
            for (auto it = parent_stack.begin(); it != parent_stack.end(); ++it) stack.push_back(it.value());
        }
        stack.push_back(to_json());
        return stack;
    }
protected:
    bool noerr_;
    std::string en, ru;
    std::string location;
    GRPCError lastError;
};

namespace grpc_errors {
    #include <gen/errors.gen.hpp>
}
