import React, { useState } from 'react';
import { contract } from '../../wailsjs/go/models';

interface RepositoryFormProps {
    onSubmit: (repo: contract.Repository) => void;
    initialData?: contract.Repository;
}

const RepositoryForm: React.FC<RepositoryFormProps> = ({ onSubmit, initialData }) => {
    const [repo, setRepo] = useState<contract.Repository>(initialData || { Owner: '', Name: '', Provider: '' });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => setRepo({ ...repo, [e.target.name]: e.target.value });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSubmit(repo);
    };

    return (
        <form onSubmit={handleSubmit} className='container'>
            <input name="Owner" placeholder="Owner" value={repo.Owner} onChange={handleChange} required />
            <input name="Name" placeholder="Name" value={repo.Name} onChange={handleChange} required />
            <input name="Provider" placeholder="Provider" value={repo.Provider} onChange={handleChange} required />
            <button type="submit">Save</button>
        </form>
    );
};

export default RepositoryForm;
