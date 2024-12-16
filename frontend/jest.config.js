const config = {
    testEnvironment: 'jsdom', // For React DOM-based testing
    transform: {
        '^.+\\.tsx?$': 'ts-jest', // Use ts-jest for TypeScript files
        '^.+\\.js$': 'babel-jest', // Use babel-jest for ES files
    },
    testMatch: ["**/__tests__/**/*.[jt]s?(x)", "**/?(*.)+(spec|test).[tj]s?(x)"],
    moduleNameMapper: {
        // Handle CSS and asset imports
        '\\.(css|scss|sass|less)$': 'identity-obj-proxy',
        '\\.(jpg|jpeg|png|gif|webp|svg)$': '<rootDir>/__mocks__/fileMock.ts',
    },
    // setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'], // Setup file for jest-dom
    testPathIgnorePatterns: ['/node_modules/', '/dist/'], // Exclude unwanted directories
};

export default config;
