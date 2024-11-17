# Assignment Submission Portal

This is going to have features for both Admin and users with OAuth2.

## For starting the server

`make build` \
`./gx api`

## User Endopints

`POST /register` - Register a new User \
`POST /login` - User login \
`POST /upload` - Upload an Assignment \
`GET /admins` - Fetches all admins

## Admin Endpoints

`POST /register` - Register a new Admin \
`POST /login` - Admin login \
`GET /assignments` - View assignments tagged to admin \
`POST /assignments/:id/accept` - Accept an assignment \
`POST /assignments/:id/reject` - Reject an assignment
