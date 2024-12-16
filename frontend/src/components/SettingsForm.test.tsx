import { cleanup, render, waitFor } from "@testing-library/react"
import SettingsForm from "./SettingsForm"


afterEach(cleanup)
describe('SettingsForm', () => {
    it("should render defaults when settings are not available", async () => {
        Object.defineProperty(window, 'go', {
            value: {
                main: {
                    App: {
                        GetSettings: () => { return {} }
                    }
                }
            }
        })
        const { getByTestId } = render(<SettingsForm onSubmit={() => { }} />)
        await waitFor(() => {
            const username = getByTestId('username') as HTMLInputElement
            const excludeDraft = getByTestId('exclude-draft') as HTMLInputElement
            const excludeConflicting = getByTestId('exclude-conflicting') as HTMLInputElement
            expect(username.value).toBe("")
            expect(excludeConflicting.checked).toBe(false)
            expect(excludeDraft.checked).toBe(false)
        })
    })

})
