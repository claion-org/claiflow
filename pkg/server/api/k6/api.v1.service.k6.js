import { check, sleep, group } from 'k6';
import { Cluster } from './api.v1.cluster.js';
import { Service } from './api.v1.service.js';
import { baseURL } from './flag.js';
import { deepCompare } from './util.js';

const cluster = new Cluster(baseURL);
const service = new Service(baseURL);

let newClusterFoo = {
    uuid: null,
    name: `service-test-foo-${Date.now()}`,
    summary: `test cluster`
}

let newClusterBar = {
    uuid: null,
    name: `service-test-bar-${Date.now()}`,
    summary: `test cluster`
}

export function setup() {
    let resultFoo = cluster.Create(newClusterFoo.uuid, newClusterFoo.name, newClusterFoo.summary);
    group(`setup: created a new cluster foo`, () => {
        check(resultFoo, { 'status is 200': (r) => r.status == 200 });
        check(JSON.parse(resultFoo.body), { 'check cluster name': (r) => r.name == newClusterFoo.name });
        check(JSON.parse(resultFoo.body), { 'check cluster summary': (r) => r.summary == newClusterFoo.summary });
    })

    let resultBar = cluster.Create(newClusterBar.uuid, newClusterBar.name, newClusterBar.summary);
    group(`setup: created a new cluster bar`, () => {
        check(resultBar, { 'status is 200': (r) => r.status == 200 });
        check(JSON.parse(resultBar.body), { 'check cluster name': (r) => r.name == newClusterBar.name });
        check(JSON.parse(resultBar.body), { 'check cluster summary': (r) => r.summary == newClusterBar.summary });
    })

    console.log(`setup: ${resultFoo.body}`)
    console.log(`setup: ${resultBar.body}`)

    return { foo: JSON.parse(resultFoo.body), bar: JSON.parse(resultBar.body) }
}


export function teardown(data) {
    console.log(`teardown: ${data.foo.uuid}`)
    console.log(`teardown: ${data.bar.uuid}`)

    group(`teardown: remove cluster foo`, () => {
        let resultFoo = cluster.Delete(data.foo.uuid);

        check(resultFoo, { 'status is 200': (r) => r.status == 200 });
    })

    group(`teardown: remove cluster bar`, () => {
        let resultBar = cluster.Delete(data.bar.uuid);

        check(resultBar, { 'status is 200': (r) => r.status == 200 });
    })
}


export default function (data) {
    group(`CreateService (foo)`, () => {
        const uuid = null;
        const name = `test-service-foo-${Date.now()}`;
        const summary = `test service`;
        const template_uuid = `example_simple`;
        const inputs = {
            name: "foo"
        };
        const subscribed_channel = null;


        let result = service.Create(data.foo.uuid, uuid, name, summary, template_uuid, inputs, subscribed_channel)

        console.log(`check inputs expected=${JSON.stringify(inputs)} actual=${JSON.stringify(JSON.parse(result.body).inputs)}`)

        check(result, { 'status is 200': (r) => r.status == 200 });
        check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
        check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
        check(JSON.parse(result.body), { 'check template_uuid': (r) => r.template_uuid == template_uuid });
        check(JSON.parse(result.body), { 'check inputs': (r) => deepCompare(r.inputs, inputs) });
        check(JSON.parse(result.body), { 'check subscribed_channel': (r) => r.subscribed_channel == subscribed_channel });

        console.log(`created a new service body=${result.body}`)
        const newService = JSON.parse(result.body)

        group(`GetService`, () => {
            let result = service.Get(newService.cluster_uuid, newService.uuid);
            console.log(`get created service body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
            check(JSON.parse(result.body), { 'check template_uuid': (r) => r.template_uuid == template_uuid });
            check(JSON.parse(result.body), { 'check inputs': (r) => deepCompare(r.inputs, inputs) });
            check(JSON.parse(result.body), { 'check subscribed_channel': (r) => r.subscribed_channel == subscribed_channel });
        })
    })

    group(`CreateService (foo, bar)`, () => {
        var clusters = [
            data.foo.uuid,
            data.bar.uuid
        ];
        const uuid = null;
        const name = `test-service-foobar-${Date.now()}`;
        const summary = `test service`;
        const template_uuid = `example_simple`;
        const inputs = {
            name: "foo"
        };
        const subscribed_channel = null;

        let result = service.CreateMultiClusters(clusters, uuid, name, summary, template_uuid, inputs, subscribed_channel)
        console.log(`created a new service body=${result.body}`)

        check(result, { 'status is 200': (r) => r.status == 200 });

        JSON.parse(result.body).forEach(e => {
            switch (e.cluster_uuid) {
                case data.foo.uuid:
                    group(`Check (foo)`, () => {
                        check(e, { 'check name': (r) => r.name == name });
                        check(e, { 'check summary': (r) => r.summary == summary });
                        check(e, { 'check template_uuid': (r) => r.template_uuid == template_uuid });
                        check(e, { 'check inputs': (r) => deepCompare(r.inputs, inputs) });
                        check(e, { 'check subscribed_channel': (r) => r.subscribed_channel == subscribed_channel });
                    })
                    break;
                case data.bar.uuid:
                    group(`Check (bar)`, () => {
                        check(e, { 'check name': (r) => r.name == name });
                        check(e, { 'check summary': (r) => r.summary == summary });
                        check(e, { 'check template_uuid': (r) => r.template_uuid == template_uuid });
                        check(e, { 'check inputs': (r) => deepCompare(r.inputs, inputs) });
                        check(e, { 'check subscribed_channel': (r) => r.subscribed_channel == subscribed_channel });
                    })
                    break;
            }
        })

        JSON.parse(result.body).forEach(e => {
            switch (e.cluster_uuid) {
                case data.foo.uuid:
                    group(`GetService (foo)`, () => {
                        let result = service.Get(e.cluster_uuid, e.uuid);
        
                        console.log(`get created service body=${result.body}`)
        
                        check(result, { 'status is 200': (r) => r.status == 200 });
                        check(e, { 'check name': (r) => r.name == name });
                        check(e, { 'check summary': (r) => r.summary == summary });
                        check(e, { 'check template_uuid': (r) => r.template_uuid == template_uuid });
                        check(e, { 'check inputs': (r) => deepCompare(r.inputs, inputs) });
                        check(e, { 'check subscribed_channel': (r) => r.subscribed_channel == subscribed_channel });
                    })
                    break;
                case data.bar.uuid:
                    group(`GetService (bar)`, () => {
                        let result = service.Get(e.cluster_uuid, e.uuid);
        
                        console.log(`get created service body=${result.body}`)
        
                        check(result, { 'status is 200': (r) => r.status == 200 });
                        check(e, { 'check name': (r) => r.name == name });
                        check(e, { 'check summary': (r) => r.summary == summary });
                        check(e, { 'check template_uuid': (r) => r.template_uuid == template_uuid });
                        check(e, { 'check inputs': (r) => deepCompare(r.inputs, inputs) });
                        check(e, { 'check subscribed_channel': (r) => r.subscribed_channel == subscribed_channel });
                    })
                    break;
            }
        })
    })
}
