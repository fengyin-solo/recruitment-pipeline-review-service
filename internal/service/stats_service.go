package service

type DepartmentStats struct {
	Department string `json:"department"`
	JobCount   int    `json:"job_count"`
}

type JobFunnelStats struct {
	JobID        string `json:"job_id"`
	JobTitle     string `json:"job_title"`
	CandidateCount int `json:"candidate_count"`
	InterviewCount int `json:"interview_count"`
	OfferCount     int `json:"offer_count"`
}

func (s *Service) DepartmentStats() ([]DepartmentStats, error) {
	jobs := s.store.ListJobs()
	m := make(map[string]int)
	for _, j := range jobs {
		m[j.Department]++
	}
	res := make([]DepartmentStats, 0, len(m))
	for dept, count := range m {
		res = append(res, DepartmentStats{Department: dept, JobCount: count})
	}
	return res, nil
}

func (s *Service) JobFunnelStats() ([]JobFunnelStats, error) {
	jobs := s.store.ListJobs()
	res := make([]JobFunnelStats, 0, len(jobs))
	for _, j := range jobs {
		st := JobFunnelStats{
			JobID:    j.ID,
			JobTitle: j.Title,
		}
		candidateIDs := make(map[string]bool)
		for _, iv := range s.store.ListInterviews() {
			if iv.JobID == j.ID {
				st.InterviewCount++
				candidateIDs[iv.CandidateID] = true
			}
		}
		for _, o := range s.store.ListOffers() {
			if o.JobID == j.ID {
				st.OfferCount++
				candidateIDs[o.CandidateID] = true
			}
		}
		st.CandidateCount = len(candidateIDs)
		res = append(res, st)
	}
	return res, nil
}
