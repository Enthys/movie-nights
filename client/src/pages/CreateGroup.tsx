import { useState } from "react"
import Group from '../types/group';
import { redirect } from "react-router-dom";

interface CreateGroupErrors {
    internal?: string
    name?: string
    description?: string
}

interface CreateGroupRequest {
    name: string;
    description: string;
}


export default function CreateGroup() {
    const [errors, setErrors] = useState({} as CreateGroupErrors);
    const [createGroupDto, setCreateGroup] = useState({ name: '', description: ''} as CreateGroupRequest);

    function createGroup(data: CreateGroupRequest): Promise<Group> {
        return new Promise<Group>((resolve) => {
            fetch("/api/groups", { method: 'POST', body: JSON.stringify(data)})
                .then((resp) => {
                    resp.json().then((data) => {
                        if (resp.status !== 201) {
                            setErrors(data.errors)
                            resolve({} as Group)
                            return
                        }

                        redirect('/groups/' + data.group.id)
                    })
                })
        })
    }

    return (<>
        <div className="container">
            {errors.internal !== undefined &&
                <div className="text-danger">{errors.internal}</div>
            }

            <form 
                className="needs-validation"
                onSubmitCapture={(e) => {
                    e.preventDefault();
                    createGroup(createGroupDto)
                }}
                noValidate={true}
            >
                <legend>Create new group</legend>
                <div className="form-group mb-3">
                    <label className="form-label">Group name</label>
                    <input
                        name="name"
                        type="text"
                        className={`form-control` + (errors.name !== undefined ? ' border-danger' : '')}
                        placeholder="Coolest group!"
                        minLength={4}
                        maxLength={24}
                        required
                        value={createGroupDto.name}
                        onChange={(event) => setCreateGroup({ ...createGroupDto, name: event.target.value })}
                    />
                    {errors.name !== undefined &&
                        <div className="text-danger">{errors.name}</div>
                    }
                </div>

                <div className="form-group mb-3">
                    <label>Group description <sub>*optional</sub></label>
                    <input
                        name="description"
                        type="text"
                        className={`form-control` + (errors.description !== undefined ? ' border-danger' : '')}
                        placeholder="(Optional)"
                        maxLength={300}
                        required
                        value={createGroupDto.description}
                        onChange={(event) => setCreateGroup({ ...createGroupDto, description: event.target.value })}
                    />
                    {errors.description !== undefined &&
                        <div className="text-danger">{errors.description}</div>
                    }
                </div>

                <button type="button" className="btn btn-lg btn-success" onClick={() => createGroup(createGroupDto)}>
                    Create
                </button>
            </form>
        </div>

    </>)
}