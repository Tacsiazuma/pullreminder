import { useEffect, useState } from 'react';
import { Link, RouterProvider, createBrowserRouter, Outlet } from 'react-router-dom';
import './App.css';
import { Repos, AddRepo, CheckPRs, UpdateSchedule, GetSettings, UpdateSettings } from "../wailsjs/go/main/App";
import { BrowserOpenURL } from '../wailsjs/runtime'
import { contract } from '../wailsjs/go/models';
import RepositoryList from './components/RepositoryList';
import RepositoryForm from './components/RepositoryForm';
import SettingsForm from './components/SettingsForm';
import ScheduleForm from './components/ScheduleForm';
import PullRequestList from './components/PullRequestList';


const Layout = () => {
    return <div>
        <nav>
            <Link to="/">Pull requests</Link> | <Link to="/repos">Repositories</Link> | <Link to="/add-repository">Add Repository</Link> |{' '}
            <Link to="/schedule">Set Schedule</Link> | <Link to="/settings">Settings</Link>
        </nav>
        <Outlet />
    </div>
}
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
    }, [])

    function handleDeleteRepository() {
    }

    function handleAddOrUpdateRepository(repo: contract.Repository) {
        AddRepo(repo).then(alert)
    }
    const router = createBrowserRouter([
        {
            element: <Layout />,
            children: [
                {
                    path: "/repos",
                    element:
                        <RepositoryList
                            repositories={repositories}
                            onDelete={handleDeleteRepository}
                            onEdit={(repo: contract.Repository) => console.log('Edit action:', repo)} // You can redirect to edit here
                        />
                },
                {
                    path: "/schedule",
                    element: <ScheduleForm onSubmit={(schedule) => UpdateSchedule(schedule)} />
                },
                {
                    path: "/",
                    element: < PullRequestList
                        pullrequests={pullrequests}
                        onVisit={BrowserOpenURL}
                    />
                },
                {
                    path: "/add-repository",
                    element: < RepositoryForm onSubmit={handleAddOrUpdateRepository} />
                },
                {
                    path: "/settings",
                    element:
                        <SettingsForm
                            onSubmit={(settings) => UpdateSettings(settings)}
                            initialSettings={{ Username: "", ExcludeDraft: false, ExcludeConflicting: false }}
                        />,
                }
            ]
        }])
    return (
        <div>
            <RouterProvider router={router}>
            </RouterProvider>
        </div>
    );
};

export default App
