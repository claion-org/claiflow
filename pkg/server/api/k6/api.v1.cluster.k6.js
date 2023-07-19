import { check, sleep, group } from 'k6';
import { Cluster } from './api.v1.cluster.js';
import { baseURL } from './flag.js';


const cluster = new Cluster(baseURL);

export default function () {
    const uuid = null
    const name = `test-cluster-${Date.now()}`
    const summary = `test cluster`
    const name_change = `${name}-changed`
    const summary_change = `${summary} changed`

    group(`create a new cluster`, () => {
        let result = cluster.Create(uuid, name, summary);
        console.log(`created a new cluster body=${result.body}`)

        check(result, { 'check status is 200': (r) => r.status == 200 });
        check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
        check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });

        let newCluster = JSON.parse(result.body)

        group(`get created cluster`, () => {
            let result = cluster.Get(newCluster.uuid);
            console.log(`get created cluster body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check name': (r) => r.name == name });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
        })

        group(`update cluster name`, () => {
            let result = cluster.Update(newCluster.uuid, name_change, null);
            console.log(`updated cluster name body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_change });
            check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });


            group(`get updated cluster`, () => {
                let result = cluster.Get(newCluster.uuid);
                console.log(`get updated cluster body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_change });
                check(JSON.parse(result.body), { 'check summary': (r) => r.summary == summary });
            })
        })

        group(`update cluster summary`, () => {
            let result = cluster.Update(newCluster.uuid,  null, summary_change);
            console.log(`updated cluster summary body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_change });
            check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_change });
       
            group(`get updated cluster`, () => {
                let result = cluster.Get(newCluster.uuid);
                console.log(`get updated cluster body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check changed name': (r) => r.name == name_change });
                check(JSON.parse(result.body), { 'check changed summary': (r) => r.summary == summary_change });
            })
        })

        group(`undo changed`, () => {
            let result = cluster.Update(newCluster.uuid,  name, summary);
            console.log(`undo changed body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });
            check(JSON.parse(result.body), { 'check undoed name': (r) => r.name == name });
            check(JSON.parse(result.body), { 'check undoed summary': (r) => r.summary == summary });
      
      
            group(`get undo cluster`, () => {
                let result = cluster.Get(newCluster.uuid);
                console.log(`get undo cluster body=${result.body}`)

                check(result, { 'status is 200': (r) => r.status == 200 });
                check(JSON.parse(result.body), { 'check undoed name': (r) => r.name == name });
                check(JSON.parse(result.body), { 'check undoed summary': (r) => r.summary == summary });
            })
        })

        group(`remove test cluster`, () => {
            let result = cluster.Delete(newCluster.uuid);
            console.log(`removed cluster body=${result.body}`)

            check(result, { 'status is 200': (r) => r.status == 200 });

            group(`get removed cluster`, () => {
                let result = cluster.Get(newCluster.uuid);
                console.log(`get removed cluster body=${result.body}`)

                check(result, { 'status is 404': (r) => r.status == 404 });
            })
        })
    })
}
