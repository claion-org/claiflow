import { check, sleep, group } from 'k6';
import { CompareDate } from './util.js';
import { Cluster } from './api.v1.cluster.js';
import { ClusterClientToken } from './api.v1.cluster_client_token.js';
import { baseURL } from './flag.js';

const cluster = new Cluster(baseURL);
const clusterClientToken = new ClusterClientToken(baseURL);

let newCluster = {
    uuid : null,
    name : `cluster-client-token-test-${Date.now()}`,
    summary: `test cluster`
}

export function setup() {
    let result = cluster.Create(newCluster.uuid, newCluster.name, newCluster.summary);

    check(result, { 'created cluster; status is 200': (r) => r.status == 200 });
    check(JSON.parse(result.body), { 'check cluster name': (r) => r.name == newCluster.name });
    check(JSON.parse(result.body), { 'check cluster summary': (r) => r.summary == newCluster.summary });

    console.log(result.body)

    return { newCluster: JSON.parse(result.body) }
}


export function teardown(data) {
    console.log(data.newCluster.uuid)
    let result = cluster.Delete(data.newCluster.uuid);

    check(result, { 'removed cluster; status is 200': (r) => r.status == 200 });
}
  
export default function (data) {
    const cluster_uuid = data.newCluster.uuid;
    const uuid        = null;
    const name        = `test-cluster-client-token-test-${Date.now()}`;
    const summary     = `test clusterClientToken`;
    const token       = null;

    // sleep(1);
    group(`create a new cluster client token`, () => {
        let result = clusterClientToken.Create(cluster_uuid, uuid, name, summary, null)
        console.log(`create a new cluster client token body=${result.body}`)

        check(result, { 'check status is 200': (r) => r.status == 200 });
        check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
        check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
        check(JSON.parse(result.body), { 'check cluster_uuid': (r) => r.cluster_uuid == cluster_uuid });

        let newClusterClentToken = JSON.parse(result.body)

        const token = newClusterClentToken.token
        const issued_at_time = new Date((new Date(newClusterClentToken.issued_at_time)).setMilliseconds(0)) // trimming ms
        const expiration_time = new Date((new Date(newClusterClentToken.expiration_time)).setMilliseconds(0)) // trimming ms
     
        console.log(`token=${token}`)
        console.log(`issued_at_time=${issued_at_time}`)
        console.log(`expiration_time=${expiration_time}`)

        const name_changed = `${newClusterClentToken.name}-changed`
        const summary_changed = `${newClusterClentToken.summary} changed`
        const token_changed = `${newClusterClentToken.token}`.toUpperCase()
        const issued_at_time_changed = new Date(issued_at_time.getTime() + 1000)
        const expiration_time_changed = new Date(expiration_time.getTime() + 1000)

        console.log(`name_changed=${name_changed}`)
        console.log(`summary_changed=${summary_changed}`)
        console.log(`token_changed=${token_changed}`)
        console.log(`issued_at_time_changed=${issued_at_time_changed} type=${typeof(newClusterClentToken.issued_at_time)}`)
        console.log(`expiration_time_changd=${expiration_time_changed}`)

        // sleep(1);
        group(`get created cluster client token`, () => {
            let result = clusterClientToken.Get(newClusterClentToken.uuid);
            console.log(`get created cluster client token body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
            check(JSON.parse(result.body), { 'check cluster_uuid': (r) => r.cluster_uuid == cluster_uuid });
        });

        // sleep(1);
        group(`update cluster client token name`, () => {
            let result = clusterClientToken.Update(newClusterClentToken.uuid, name_changed, null, null, null, null);
            console.log(`updated cluster client token name body=${result.body}`)

            console.log(`${JSON.parse(result.body).issued_at_time} === ${issued_at_time} CompareDate=${CompareDate(JSON.parse(result.body).issued_at_time, issued_at_time)}`)
            console.log(`${JSON.parse(result.body).expiration_time} === ${expiration_time} CompareDate=${CompareDate(JSON.parse(result.body).expiration_time, expiration_time)}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
            check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary });
            check(JSON.parse(result.body), { 'check token': (r) => r.token == token });
            check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time) });
            check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
       
            // sleep(1);
            group(`get updated cluster client token`, () => {
                let result = clusterClientToken.Get(newClusterClentToken.uuid);
                console.log(`get updated cluster client token body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
                check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary });
                check(JSON.parse(result.body), { 'check token': (r) => r.token == token });
                check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time) });
                check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
            });
        });

        // sleep(1);
        group(`update cluster client token summary`, () => {
            let result = clusterClientToken.Update(newClusterClentToken.uuid, null, summary_changed, null, null, null);
            console.log(`updated cluster client token summary body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
            check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_changed });
            check(JSON.parse(result.body), { 'check token': (r) => r.token == token });
            check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time) });
            check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
       
            // sleep(1);
            group(`get updated cluster client token`, () => {
                let result = clusterClientToken.Get(newClusterClentToken.uuid);
                console.log(`get updated cluster client token body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
                check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_changed });
                check(JSON.parse(result.body), { 'check token': (r) => r.token == token });
                check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time) });
                check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
            });
        });
        
        // sleep(1);
        group(`update cluster client token token`, () => {
            let result = clusterClientToken.Update(newClusterClentToken.uuid, null, null, token_changed, null, null);
            console.log(`updated cluster client token token body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
            check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_changed });
            check(JSON.parse(result.body), { 'check token': (r) => r.token == token_changed });
            check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time) });
            check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
       
            // sleep(1);
            group(`get updated cluster client token`, () => {
                let result = clusterClientToken.Get(newClusterClentToken.uuid);
                console.log(`get updated cluster client token body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
                check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_changed });
                check(JSON.parse(result.body), { 'check token': (r) => r.token == token_changed });
                check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time) });
                check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
            });
        });

        // sleep(1);
        group(`update cluster client token issued_at_time`, () => {
            let result = clusterClientToken.Update(newClusterClentToken.uuid, null, null, null, issued_at_time_changed, null);
            console.log(`updated cluster client token issued_at_time body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
            check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_changed });
            check(JSON.parse(result.body), { 'check token': (r) => r.token == token_changed });
            check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time_changed) });
            check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
       
            // sleep(1);
            group(`get updated cluster client token`, () => {
                let result = clusterClientToken.Get(newClusterClentToken.uuid);
                console.log(`get updated cluster client token body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
                check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_changed });
                check(JSON.parse(result.body), { 'check token': (r) => r.token == token_changed });
                check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time_changed) });
                check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
            });
        });

        // sleep(1);
        group(`update cluster client token expiration_time`, () => {
            let result = clusterClientToken.Update(newClusterClentToken.uuid, null, null, null, null, expiration_time_changed);
            console.log(`updated cluster expiration_time body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
            check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_changed });
            check(JSON.parse(result.body), { 'check token': (r) => r.token == token_changed });
            check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time_changed) });
            check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time_changed) });
       
            // sleep(1);
            group(`get updated cluster client token`, () => {
                let result = clusterClientToken.Get(newClusterClentToken.uuid);
                console.log(`get updated cluster client token body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_changed });
                check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_changed });
                check(JSON.parse(result.body), { 'check token': (r) => r.token == token_changed });
                check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time_changed) });
                check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time_changed) });
            });
        });

        // sleep(1);
        group(`undo udated cluster client token properties`, () => {
            let result = clusterClientToken.Update(newClusterClentToken.uuid, name, summary, token, issued_at_time, expiration_time);
            console.log(`updo updated cluster client token properties body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name });
            check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary });
            check(JSON.parse(result.body), { 'check token': (r) => r.token == token });
            check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time) });
            check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
       
            // sleep(1);
            group(`get updated cluster client token`, () => {
                let result = clusterClientToken.Get(newClusterClentToken.uuid);
                console.log(`get updated cluster client token body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name });
                check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary });
                check(JSON.parse(result.body), { 'check token': (r) => r.token == token });
                check(JSON.parse(result.body), { 'check issued_at_time': (r) => 0 == CompareDate(r.issued_at_time, issued_at_time) });
                check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 == CompareDate(r.expiration_time, expiration_time) });
            });
        });

        // sleep(1);
        group(`refresh cluster client token properties`, () => {
            let result = clusterClientToken.Refresh(newClusterClentToken.uuid);
            console.log(`refreshed cluster client token properties body=${result.body}`)

            console.log(`refreshed cluster client token: ${JSON.parse(result.body).expiration_time} === ${expiration_time} CompareDate=${CompareDate(JSON.parse(result.body).expiration_time, expiration_time)}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 > CompareDate(r.expiration_time, expiration_time_changed) });
        });

        // sleep(1);
        group(`expire cluster client token properties`, () => {
            let result = clusterClientToken.Expire(newClusterClentToken.uuid);
            console.log(`expired cluster client token properties body=${result.body}`)

            console.log(`expired cluster client token: ${JSON.parse(result.body).expiration_time} === ${new Date()} CompareDate=${CompareDate(JSON.parse(result.body).expiration_time, new Date())}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check expiration_time': (r) => 0 >= CompareDate(r.expiration_time, new Date()) });
        });

        // sleep(1);
        group(`remove test cluster client token`, () => {
            let result = clusterClientToken.Delete(newClusterClentToken.uuid);
            console.log(`removed test cluster client token body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            
            // sleep(1);
            group(`get removed cluster client token`, () => {
                let result = clusterClientToken.Get(newClusterClentToken.uuid);
                console.log(`get removed cluster client token body=${result.body}`)

                check(result, { 'status is 404': (r) => r.status == 404 });
            })
        });
    })
}
