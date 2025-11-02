stack:
GO,Nuxt,MongoDB
Run:

npm install
npm run dev

Ensure backend is running at http://localhost:8080 (main.go)

env:
1.backend
- ./
2.frontend
- ./
- nuxt.config.ts
- frontend/app/components/cardmodal.vue const props = defineProps({