import { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import './App.css';
import { Repos, AddRepo, CheckPRs, UpdateSchedule } from "../wailsjs/go/main/App";
import { BrowserOpenURL } from '../wailsjs/runtime'
import { contract } from '../wailsjs/go/models';
import RepositoryList from './components/RepositoryList';
import RepositoryForm from './components/RepositoryForm';
import SettingsForm from './components/SettingsForm';
import ScheduleForm from './components/ScheduleForm';
import PullRequestList from './components/PullRequestList';

function App() {
    const [repositories, setRepos] = useState(Array<contract.Repository>);
    const [pullrequests, setPRs] = useState(Array<contract.Pullrequest>);

    useEffect(() => {
        const fetchRepos = async () => {
            const repos = await Repos()
            setRepos(repos)
        }
        const fetchPRs = async () => {
            const prs = await CheckPRs()
            setPRs(prs)
        }
        fetchRepos()
        fetchPRs()
    })

    function handleDeleteRepository() {
    }

    function handleAddOrUpdateRepository(repo: contract.Repository) {
        AddRepo(repo).then(alert)
    }

    return (
        <Router>
            <nav>
                <Link to="/">Pull requests</Link> | <Link to="/repos">Repositories</Link> | <Link to="/add-repository">Add Repository</Link> |{' '}
                <Link to="/schedule">Set Schedule</Link> | <Link to="/settings">Settings</Link>
            </nav>
            <Routes>
                <Route
                    path="/repos"
                    element={
                        <RepositoryList
                            repositories={repositories}
                            onDelete={handleDeleteRepository}
                            onEdit={(repo: contract.Repository) => console.log('Edit action:', repo)} // You can redirect to edit here
                        />
                    }
                />
                <Route
                    path="/"
                    element={
                        <PullRequestList
                            pullrequests={pullrequests}
                            onVisit={BrowserOpenURL}
                        />
                    }
                />
                <Route
                    path="/add-repository"
                    element={<RepositoryForm onSubmit={handleAddOrUpdateRepository} />}
                />
                <Route
                    path="/schedule"
                    element={<ScheduleForm onSubmit={(schedule) => UpdateSchedule(schedule)} />}
                />
                <Route
                    path="/settings"
                    element={
                        <SettingsForm
                            onSubmit={(settings) => console.log('Save Settings:', settings)}
                            initialSettings={{ username: '', includeConflicting: true, includeDraft: false }}
                        />
                    }
                />
            </Routes>
        </Router>
    );
};

export default App
