import { useState } from 'react';
import { default as IMovie} from '../types/movie';
import addMovie from '../functions/addMovie';

interface MovieAddArgs {
    onAdd: (movie: IMovie) => void
    onFail: (err: string) => void
}

export default function MovieAdd({ onAdd, onFail }: MovieAddArgs) {
    const [link, setLink] = useState('')
    const [err, setError] = useState('')

    return (
        <div className="input-group" hx-ext="response-targets">
            <input 
                name="url" 
                type="text" 
                className="form-control" 
                placeholder="Movie IMDb link"
                value={link}
                onChange={(event) => {
                    try {
                        setLink(event.target.value)
                        setError('')
                        new URL(event.target.value)
                    } catch (err) {
                        setError('invalid link')
                    }
                }}
            />

            <button 
                type="button"
                className="btn btn-outline-success"
                disabled={err != ''}
                onClick={() => addMovie(link).then(onAdd).catch(onFail)}
            >Add</button>
        </div>
    )
}