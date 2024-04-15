import {useState} from 'react';
import Group from '../components/Group';
import { default as IGroup} from '../types/group';

function getGroups(nameSearch: string): Promise<IGroup[]> {
    return new Promise<IGroup[]>((resolve, reject) => {
        fetch("/api/groups"+ (nameSearch !== '' ? `/search?name=${nameSearch}` : ''))
            .then((resp) => {
                if (resp.status !== 200) {
                    reject(new Error("failed to retrieve groups"))
                }

                return resp.json();
            })
            .then((data) => {
                resolve(data.groups);
            })
    })
}

export default function Groups() {
    const [isSet, setIsSet] = useState(false)
    const [groups, setGroups] = useState([] as IGroup[]);
    const [search, setSearch] = useState('')

    if (!isSet) {
        getGroups('').then((groups) => {
            setGroups(groups);
            setIsSet(true);
        });
    }

    return (
        <>
            <div className="container">
                <div className="row">
                    <form onSubmit={() => getGroups(search).then(setGroups)} className="col-10">
                        <div className="input-group">
                            <input
                                className="form-control"
                                id="group-name-search"
                                name="name"
                                type="text"
                                value={search}
                                onChange={(event) => {
                                    setSearch(event.target.value);
                                }}
                                placeholder="Search for group"
                            />

                            <button 
                                type="button" 
                                className="btn btn-outline-secondary"
                                onClick={() => {
                                    getGroups(search).then(setGroups);
                                }}
                            >Search</button>
                        </div>
                    </form>

                    <div className="col-2">
                        <a href="/groups/create"><button className="btn btn-success">Create group</button></a>
                    </div>
                </div>

                {groups.length === 0 && 
                    <div className="row text-center">
                        <h4>You are not a part of any group.</h4>
                    </div>
                }

                <div id="groups-collection">
                    {groups.map((group) => <Group group={group} />)}
                </div>
            </div>
            
        </>
    )
}
