import React from 'react';
import { contract } from '../../wailsjs/go/models';

interface RepositoryListProps {
    repositories: contract.Repository[];
    onDelete: (repo: contract.Repository) => void;
    onEdit: (repo: contract.Repository) => void;
}
const RepositoryList: React.FC<RepositoryListProps> = ({ repositories,  onDelete, onEdit }) => (
    <div className='container'> <h2>Repositories</h2>
        <table>
            <thead>
                <tr>
                    <th>Owner</th>
                    <th>Name</th>
                    <th>Provider</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {repositories.map((repo) => (
                    <tr key={repo.Owner + repo.Name}>
                        <td>{repo.Owner}</td>
                        <td>{repo.Name}</td>
                        <td>{repo.Provider}</td>
                        <td>
                            <button onClick={() => onEdit(repo)}>Edit</button>
                            <button onClick={() => onDelete(repo)}>Delete</button>
                        </td>
                    </tr>
                ))}
            </tbody>
        </table>
    </div>
);

export default RepositoryList;
