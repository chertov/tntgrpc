#include <memory>
#include <iostream>
#include <string>
#include <thread>
#include <unistd.h>

#include <grpcpp/grpcpp.h>
#include <grpc/support/log.h>
#include "date.h"

#include "luatools.hpp"
#include "errors.hpp"

#ifdef TARANTOOL_GRPC
    #include <tarantool/module.h>
#endif

enum CallStatus { PROCESS, FINISH };
class CallMethod {
public:
    int fd[2];
    GRPCError err = NOERR;
    CallMethod() {
        int r = pipe(fd);
        if (r == -1) {
            std::cerr << strerror (errno) << "\n";
        }
    }
    virtual ~CallMethod() {
        close(fd[0]);
        close(fd[1]);
    }
    virtual void Proceed() = 0;
    virtual void Run() = 0;
    void wait() {
#ifdef TARANTOOL_GRPC
        int release = 111111;
        // { std::stringstream ss; ss << std::chrono::system_clock::now() << "    " << "    " << "wait for release" << std::endl; std::cout << ss.str() << std::flush; }
        read(fd[0], &release, sizeof(int));
        // { std::stringstream ss; ss << std::chrono::system_clock::now() << "    " << "    " << "lock was released = " << release << std::endl; std::cout << ss.str() << std::flush; }
#endif
    }
};

struct PipeData{
    CallMethod *ptr = nullptr;
    int fd;
};

int pipefd[2];
void Runner(CallMethod *ptr) {
#ifdef TARANTOOL_GRPC
    PipeData pd;
    pd.ptr = ptr;
    write(pipefd[1], &pd, sizeof(PipeData));
    // std::cout << "write " << "   " << pd.fd << "    " << pd.ptr << std::endl;
#else
    ptr->Run();
#endif
}


int CallbackSaveContext(lua_State* L) {
    int argc = lua_gettop(L);
    if (argc != 1) {
        luaL_error(L, "no argument was given! Do you used this function like 'context.Save()'? Replace '.' to ':'.\nYou must use this function like 'context:Save()'");
        return 0;
    }
    if (lua_type(L, 1) != LUA_TTABLE) {
        luaL_error(L, "first argument isn't table!\nYou must use this function like 'context:Save()'");
        return 0;
    }
    grpc::ServerContext *context = nullptr;
    lua_getfield (L, 1, "_ptr");
    if (lua_type(L, -1) != LUA_TLIGHTUSERDATA) {
        luaL_error(L, "first argument doesn't have key '_ptr' with userdata type!\nYou must use this function like 'context:Save()'");
        return 0;
    }
    context = reinterpret_cast<grpc::ServerContext *>(lua_touserdata(L, -1));
//    if (context) {
//        std::cout << "Save context!" << std::endl;
//    }
    return 0;
}

GRPCError tolua_context(lua_State *L, grpc::ServerContext* context) {
    lua_newtable(L);
    lua_pushstring(L, "_ptr");
    lua_pushlightuserdata(L, context);
    lua_settable(L, -3);
    {
        lua_pushstring(L, "Save");
        lua_pushcfunction(L, CallbackSaveContext);
        lua_settable(L, -3);
    }
    {
        lua_pushstring(L, "client_metadata");
        lua_newtable(L);
        auto client_metadata = context->client_metadata();
        for (auto it=client_metadata.begin(); it!=client_metadata.end(); ++it) {
            auto key = (*it).first;
            auto value = (*it).second;
            lua_pushlstring(L, key.data(), key.size());
            lua_pushlstring(L, value.data(), value.size());
            lua_settable(L, -3);
        }
        lua_settable(L, -3);
    }
    {
        auto peer = context->peer();
        lua_pushstring(L, "peer");
        lua_pushlstring(L, peer.data(), peer.size());
        lua_settable(L, -3);
    }
    {
        lua_pushstring(L, "auth_context");
        lua_newtable(L);
        auto auth_context = context->auth_context();
        for (auto it=auth_context->begin(); it!=auth_context->end(); ++it) {
            auto key = (*it).first;
            auto value = (*it).second;
            lua_pushlstring(L, key.data(), key.size());
            lua_pushlstring(L, value.data(), value.size());
            lua_settable(L, -3);
        }
        lua_settable(L, -3);
    }
    return NOERR;
}

bool fieldNameCompare(const std::string &name1, const std::string &name2) {
    if (name1.size() != name2.size()) return false;
    for (std::string::const_iterator c1 = name1.begin(), c2 = name2.begin(); c1 != name1.end(); ++c1, ++c2) {
        if (tolower(*c1) != tolower(*c2)) return false;
    }
    return true;
}

#include "gen/grpc.gen.cpp"

using namespace date;

std::unique_ptr<grpc::Server> server;
std::unique_ptr<grpc::ServerCompletionQueue> cq;

bool run_server(const std::string server_address, bool sync_mode) {
    if (sync_mode) {
        grpc::ServerBuilder builder;
        builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
        RegisterServicesSync(builder);
        std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
        { std::stringstream ss; ss << std::chrono::system_clock::now() << "    " << "Server listening on " << server_address << " in sync mode" << std::endl; std::cout << ss.str() << std::flush; }
        server->Wait();
    } else {
        grpc::ServerBuilder builder;
        builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
        RegisterServicesAsync(builder);
        cq = builder.AddCompletionQueue();
        server = builder.BuildAndStart();
        InitRPCsAsync(cq);
        { std::stringstream ss; std::cout << std::chrono::system_clock::now() << "    " << "Server listening on " << server_address << " in async mode" << std::endl; std::cout << ss.str() << std::flush; }
        while (true) {
            void* tag;
            bool ok = false;
            // auto r = cq->AsyncNext(&tag, &ok, gpr_time_0(GPR_CLOCK_REALTIME));
            auto r = cq->Next(&tag, &ok);
            if (r == grpc::CompletionQueue::TIMEOUT) {
            } else if (r == grpc::CompletionQueue::GOT_EVENT) {
                GPR_ASSERT(ok);
                static_cast<CallMethod*>(tag)->Proceed();
            } else if (r == grpc::CompletionQueue::SHUTDOWN ) {
                std::cout << "SHUTDOWN " << ok << std::endl;
            }
        }
    }
}

void grpc_read() {
#ifdef TARANTOOL_GRPC
    while( coio_wait(pipefd[0], COIO_READ, 365*86400*100.0) ) {
        PipeData pd;
        int r = read(pipefd[0], &pd, sizeof(PipeData));
        // std::cout << "read " << r << "   " << pd.fd << "    " << pd.ptr << std::endl;
        if (pd.ptr) {
            pd.ptr->Run();
            static int r = 10;
            write(pd.ptr->fd[1], &r, sizeof(int));
        }
    }
#endif
}

bool grpcSync = false;
std::string server_address;
void *pipe_f(void *arg) {
    run_server(server_address, grpcSync);
}
void grpc_start(const std::string address, bool sync) {
    server_address = address;
    grpcSync = sync;
    int r = pipe(pipefd);
    if (r == -1) {
        std::cerr << strerror (errno) << "\n";
    }
    pthread_t thread;
    pthread_create(&thread, nullptr, pipe_f, nullptr);
    grpc_read();
}
