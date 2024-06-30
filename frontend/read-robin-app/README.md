# Read Robin Frontend

## Overview

This project is the frontend part of the Read Robin application. It is a React application integrated with Firebase for authentication and hosting.

## Local Development

To develop locally within the devcontainer, follow these steps:

1. **Build the Application**:

   ```sh
   npm run build
   ```

2. **Start Firebase Emulators**:
   ```sh
   firebase emulators:start
   ```

_There should be a way to start the app with `npm run start` in the container, but it has not been configured yet._

## Deploying Changes to Firebase

To deploy changes to Firebase, you have two options:

1. **Merge to `main` Branch**:

   - Merging your changes to the `main` branch will automatically trigger a deployment.

2. **Manual Deployment**:
   - Build the application:
     ```sh
     npm run build
     ```
   - Deploy to Firebase:
     ```sh
     firebase deploy
     ```

## Common Issues

- Ensure you have a `firebase.js` file in the `src` directory. This file should contain your Firebase configuration and initialization.

## Styling with Material-UI

To enhance the styling of your application using Material-UI, install the necessary packages:

```sh
npm install @mui/material @emotion/react @emotion/styled
```

## Emergency Procedures

In case of an emergency where you need to pause hosting, deploy a maintenance page:

1. Create a maintenance.html file with your maintenance message.
2. Deploy the maintenance.html page to Firebase Hosting.
