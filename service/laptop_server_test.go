package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/idirall22/grpc/sample"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/idirall22/grpc/pb"
	"github.com/idirall22/grpc/service"
)

func TestCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewImageStore("../tmp")
	server := service.NewLaptopServer(laptopStore, imageStore)

	laptop := sample.NewLaptop()

	laptopNoID := sample.NewLaptop()
	laptopNoID.Id = ""

	laptopInvalidID := sample.NewLaptop()
	laptopInvalidID.Id = "invalid"

	testCases := []struct {
		name   string
		laptop *pb.Laptop
		code   codes.Code
	}{
		{
			name:   "success_with_id",
			laptop: laptop,
			code:   codes.OK,
		},
		{
			name:   "success_without_id",
			laptop: laptopNoID,
			code:   codes.OK,
		},
		{
			name:   "error_with_invalid_id",
			laptop: laptopInvalidID,
			code:   codes.InvalidArgument,
		},
		{
			name:   "error_duplicate_id",
			laptop: laptop,
			code:   codes.AlreadyExists,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			res, err := server.CreateLaptop(ctx, &pb.CreateLaptopRequest{Laptop: tc.laptop})

			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
				if len(tc.laptop.Id) > 0 {
					require.Equal(t, tc.laptop.Id, res.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}

		})
	}
}
