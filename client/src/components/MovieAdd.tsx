import { useState } from 'react';
import { default as IMovie} from '../types/movie';
import addMovie from '../functions/addMovie';

interface MovieAddArgs {
    onAdd: (movie: IMovie) => void
    onFail: (err: string) => void
}

export default function MovieAdd({ onAdd, onFail }: MovieAddArgs) {
    const [link, setLink] = useState('')

    return (
        <div className="input-group" hx-ext="response-targets">
            <input 
                name="url" 
                type="text" 
                className="form-control" 
                placeholder="Movie IMDb link"
                value={link}
                onChange={(event) => setLink(event.target.value)}
            />
            <button 
                type="button"
                className="btn btn-outline-success"
                onClick={() => addMovie(link).then(onAdd).catch(onFail)}
            >Add</button>
        </div>
    )
}