package main

import (
    "github.com/emicklei/proto"
    "sort"
)

type Service struct {
    *proto.Service
    Package string
    RPCs RPCsMap
}
type ServicesMap map[string]Service;
func (h ServicesMap)getSortedKeys() []string {
    names := make([]string, 0, len(h))
    for name := range h {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}
type RPCsMap map[string]*proto.RPC;
func (h RPCsMap)getSortedKeys() []string {
    names := make([]string, 0, len(h))
    for name := range h {
        names = append(names, name)
    }
    sort.Strings(names)
    return names
}

func NewService(ptr *proto.Service, pkgName string, pd *ProtoData, pad string) {
    service := Service{}
    service.Package = pkgName
    service.Service = ptr
    service.RPCs = make(map[string]*proto.RPC)
    log_(pad, "Service ", service.Name)
    for _, element := range ptr.Elements {
        rpc, ok := element.(*proto.RPC)
        if ok {
            log_(pad + "    ", "RPC ", rpc.Name)
            service.RPCs[pkgName + "." + service.Name + "." + rpc.Name] = rpc
        }
    }
    pd.Services[pkgName + "." + service.Name] = service
}


func (h ServicesMap)genCpp() string {
    cpp := ""
    cpp += h.genAsyncCpp()
    cpp += "\n\n"
    cpp += h.genSyncCpp()
    return cpp
}

func (h ServicesMap)genAsyncCpp() string {
    cpp := ""

    classes := make([]CppClass, 0)
    for _, path := range h.getSortedKeys() {
        service := h[path]
        cpp += "// classes for async GRPC service\n"
        cpp += "namespace " + service.Package + " {\n"
        srvname := service.Name
        cpp += "    " + srvname + "::AsyncService " + srvname + "_service;\n"
        for _, method := range service.RPCs.getSortedKeys() {
            rpc := service.RPCs[method]
            className := CppClass{service.Package,  srvname + "_" + rpc.Name}
            classes = append(classes, className)
            cpp += "    " + "class " + className.Name + " final : public CallMethod {\n"
            cpp += "    " + "public:\n"
            cpp += "        " + className.Name + "(std::unique_ptr<grpc::ServerCompletionQueue> &cq) : cq(cq), responder(&ctx) {\n"
            cpp += "            " + srvname + "_service.Request" + rpc.Name + "(&ctx, &request, &responder, cq.get(), cq.get(), this);\n"
            cpp += "        " + "}\n"
            cpp += "        " + "void Run() {\n"
            cpp += "            " + "lua_State *L = getLuaState();\n"
            cpp += "            " + "err = loadProc(L, \"" + srvname + "\",\"" + rpc.Name + "\");\n"
            cpp += "            " + "if (err->noerr()) {\n"
            cpp += "            " + "    tolua_context(L, &ctx);\n"
            cpp += "            " + "    tolua(L, \"\", request);\n"
            cpp += "            " + "    err = callProc(L, \"" + srvname + "\",\"" + rpc.Name + "\");\n"
            cpp += "            " + "    if (err->noerr()) {\n"
            cpp += "            " + "        lua_pop(L, 1); // error return value\n"
            cpp += "            " + "        err = fromlua(L, \"\", &reply);\n"
            cpp += "            " + "        if (err->noerr()) {\n"
            cpp += "            " + "            lua_pop(L, 1); // function\n"
            cpp += "            " + "            lua_pop(L, 1); // Service\n"
            cpp += "            " + "        };\n"
            cpp += "            " + "    };\n"
            cpp += "            " + "};\n"
            cpp += "        " + "}\n"
            cpp += "        " + "void Proceed() {\n"
            cpp += "            " + "if (status == PROCESS) {\n"
            cpp += "                " + "new " + className.Name + "(cq);\n"
            cpp += "                " + "Runner(this);\n"
            cpp += "                " + "wait();\n"
            cpp += "            " + "    if(!err->noerr()) {\n"
            cpp += "            " + "        std::string str = err->json();\n"
            cpp += "            " + "        std::cerr << str << std::endl;\n"
            cpp += "            " + "        responder.FinishWithError(grpc::Status(grpc::StatusCode::INTERNAL, err->json()), this);\n"
            cpp += "            " + "    } else responder.Finish(reply, grpc::Status::OK, this);\n"
            cpp += "            " + "    status = FINISH;\n"
            cpp += "            " + "} else {\n"
            cpp += "                " + "GPR_ASSERT(status == FINISH);\n"
            cpp += "                " + "delete this;\n"
            cpp += "            " + "}\n"
            cpp += "        " + "}\n"
            cpp += "    private:\n"
            cpp += "        " + rpc.ReturnsType + " reply;\n"
            cpp += "        " + rpc.RequestType + " request;\n"
            cpp += "        " + "grpc::ServerAsyncResponseWriter<" + rpc.ReturnsType + "> responder;\n"
            cpp += "        " + "grpc::ServerContext ctx;\n"
            cpp += "        " + "CallStatus status = PROCESS;\n"
            cpp += "        " + "std::unique_ptr<grpc::ServerCompletionQueue> &cq;\n"
            cpp += "    };\n"
        }
        cpp += "\n"
        cpp += "}\n\n"
    }
    cpp += "void InitRPCsAsync(std::unique_ptr<grpc::ServerCompletionQueue> &cq) {\n"
    for _, className := range classes {
        cpp += "    new " + className.NameSpace + "::" + className.Name + "(cq);\n"
    }
    cpp += "}\n\n"
    cpp += "void RegisterServicesAsync(grpc::ServerBuilder &builder) {\n"
    for _, path := range h.getSortedKeys() {
        service := h[path]
        cpp += "    builder.RegisterService(&" + service.Package + "::" +service.Name + "_service);\n"
    }
    cpp += "}\n"

    return cpp
}

func (h ServicesMap)genSyncCpp() string {
    cpp := ""
    for _, path := range h.getSortedKeys() {
        service := h[path]
        cpp += "// classes for sync GRPC service\n"
        cpp += "namespace " + service.Package + " {\n"

        for _, method := range service.RPCs.getSortedKeys() {
            rpc := service.RPCs[method]
            className := service.Name + "_" + rpc.Name + "_sync"
            cpp_srv := "class " + className + " final : public CallMethod {\n"
            cpp_srv += "public:\n"
            cpp_srv += "    " + className + "(grpc::ServerContext* context, const " + rpc.RequestType + " *request, " + rpc.ReturnsType + " *reply) : context(context), request(request), reply(reply) {}\n"
            cpp_srv += "    void Proceed() override {}\n"
            cpp_srv += "    void Run() override {\n"
            cpp_srv += "        lua_State *L = getLuaState();\n"
            cpp_srv += "        err = loadProc(L, \"" + service.Name + "\", \"" + rpc.Name + "\");\n"
            cpp_srv += "        if(err->noerr()) {\n"
            cpp_srv += "            tolua_context(L, context);\n"
            cpp_srv += "            tolua(L, \"\", *request);\n"
            cpp_srv += "            err = callProc(L, \"" + service.Name + "\", \"" + rpc.Name + "\");\n"
            cpp_srv += "            if (err->noerr()) {\n"
            cpp_srv += "                lua_pop(L, 1);                  // error return value\n"
            cpp_srv += "                err = fromlua(L, \"\", reply);    // result\n"
            cpp_srv += "                if (err->noerr()) {\n"
            cpp_srv += "                    lua_pop(L, 1);                  // function\n"
            cpp_srv += "                    lua_pop(L, 1);                  // service\n"
            cpp_srv += "                }\n"
            cpp_srv += "            }\n"
            cpp_srv += "        }\n"
            cpp_srv += "    }\n"
            cpp_srv += "private:\n"
            cpp_srv += "    " + "const " + rpc.RequestType + " *request;\n"
            cpp_srv += "    " + rpc.ReturnsType + " *reply;\n"
            cpp_srv += "    " + "grpc::ServerContext* context;\n"
            cpp_srv += "};\n"
            cpp += pad(1, cpp_srv)
        }

        cpp_srv := ""
        cpp_srv += "class " + service.Name + "ServiceSyncImpl final : public " + service.Package + "::" + service.Name + "::Service {\n"
        for _, method := range service.RPCs.getSortedKeys() {
            rpc := service.RPCs[method]
            cpp_rpc := "grpc::Status " + rpc.Name + "(grpc::ServerContext* context, const " + rpc.RequestType + " *request, " + rpc.ReturnsType + " *reply) override {\n"
            cpp_rpc += "    " + service.Name + "_" + rpc.Name + "_sync call(context,request,reply); Runner(&call); call.wait();\n"
            cpp_rpc += "    if(call.err->noerr()) { return grpc::Status::OK; } \n"
            // cpp_rpc += "    return grpc::Status(call.err->errorCode, call.err->str()); \n"
            cpp_rpc += "    return grpc::Status(grpc::StatusCode::INTERNAL, call.err->json()); \n"
            cpp_rpc += "};\n"
            cpp_srv += pad(1, cpp_rpc)
        }
        cpp_srv += "};\n"
        cpp_srv += service.Name + "ServiceSyncImpl " + service.Name + "ServiceSyncImpl_;\n"
        cpp += pad(1, cpp_srv)
        cpp += "}\n"
    }

    cpp += "void RegisterServicesSync(grpc::ServerBuilder &builder) {\n"
    for _, path := range h.getSortedKeys() {
        service := h[path]
        cpp += "    builder.RegisterService(&" + service.Package + "::" + service.Name + "ServiceSyncImpl_);\n"
    }
    cpp += "}\n"
    return cpp
}
