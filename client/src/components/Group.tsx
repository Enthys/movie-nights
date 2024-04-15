import { default as IGroup} from '../types/group';

type GroupProps = {group: IGroup}

export default function Group({group}: GroupProps) {
    return (
        <>
            <div className="card mt-3 p-3">
                <div className="row g-0">
                    <img
                        className="col-md-3 img-fluid rounded-start"
                        style={{maxHeight: "200px", objectFit: 'cover'}}
                        src={`https://api.dicebear.com/8.x/bottts/svg?seed=${group.name}`}
                    />

                    <div className="col-md-9">
                        <div className="card-body">
                            <h5 className="card-title">{group.name}</h5>
                            <p className="card-text">{group.description}</p>
                            <a href={`/groups/${group.id}`} className="btn btn-success stretched-link">Open</a>
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}
