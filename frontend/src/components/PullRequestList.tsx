import React from 'react';
import { contract } from '../../wailsjs/go/models';

interface PullRequestListProps {
    pullrequests: contract.Pullrequest[];
    onVisit: (url: string) => void;
}
const PullRequestList: React.FC<PullRequestListProps> = ({ pullrequests, onVisit }) => (
    <div className='container'>
        <h2>Pull Requests</h2>
        <table>
            <thead>
                <tr>
                    <th>Number</th>
                    <th>Title</th>
                    <th>Author</th>
                    <th>Opened</th>
                    <th>Reviewers</th>
                    <th>Mergeable</th>
                    <th>URL</th>
                </tr>
            </thead>
            <tbody>
                {pullrequests.map((pr) => (
                    <tr key={pr.Number}>
                        <td>#{pr.Number}</td>
                        <td>{pr.Title}</td>
                        <td>{pr.Author}</td>
                        <td>{pr.Opened}</td>
                        <td>{pr.Reviewers.join(', ')}</td>
                        <td>{pr.Mergeable ? "true" : "false"}</td>
                        <td><button onClick={() => onVisit(pr.URL)}>Visit</button></td>
                    </tr>
                ))}
            </tbody>
        </table>
    </div>
);

export default PullRequestList;
