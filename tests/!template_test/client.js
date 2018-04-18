const grpc = require('grpc');
process.exitCode = 1; // Exit with error if not to call process.exit(0) manually!

const tests = grpc.load(__dirname + '/tests.proto').tests;
const testService = new tests.TestService('localhost:50051', grpc.credentials.createInsecure());

/*
// very simple request test
const msg = 'test_message';
testService.Procedure({message: msg}, (err, reply) => {
    if (err) {
        console.log('err:', err);
        return;
    }
    console.log('reply:', reply);
    if (reply.message === msg) process.exit(0);
});
*/

// complex multiply requests
const promises = [];
for (let i = 0; i < 10; i++) {
    let p;
    p = new Promise((resolve, reject) => {
        const msg = 'test_message';
        testService.Procedure({message: msg}, (err, reply) => {
            if (!err && reply && reply.message === msg) { resolve(); return; }
            reject('testService.Procedure failed');
        });
    });
    promises.push(p);

    p = new Promise((resolve, reject) => {
        testService.DoesNotExistsProcedure({}, (err, reply) => {
            if (err) {
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'ProcedureNotFound') { resolve(); return; }
            }
            reject('testService.DoesNotExistsProcedure failed');
        })
    });
    promises.push(p);

    p = new Promise((resolve, reject) => {
        testService.ExistsProcedure({message: 'm1'}, (err, reply) => {
            if (!err && reply && reply.message === 'm1') { resolve(); return; }
            reject('testService2.ExistsProcedure failed');
        });
    });
    promises.push(p);

    p = new Promise((resolve, reject) => {
        testService.ErrorWhileCallProcedure({}, (err, reply) => {
            if (err) {
                // console.log('err.details', err.details);
                const grpcerr = JSON.parse(err.details);
                if (grpcerr.type === 'CallProcedureError') { resolve(); return; }
            }
            reject('testService.ErrorWhileCallProcedure failed');
        })
    });
    promises.push(p);
}

Promise.all(promises)
    .then(() => process.exit(0))
    .catch((ex) => {
        console.log(ex);
    });
